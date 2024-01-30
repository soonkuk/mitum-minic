package digest

import (
	"context"
	"sort"
	"sync"
	"time"

	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	isaacblock "github.com/ProtoconNet/mitum2/isaac/block"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/fixedtree"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Digester struct {
	sync.RWMutex
	*util.ContextDaemon
	*logging.Logging
	database      *currencydigest.Database
	localfsRoot   string
	blockChan     chan base.BlockMap
	errChan       chan error
	sourceReaders *isaac.BlockItemReaders
	fromRemotes   isaac.RemotesBlockItemReadFunc
	networkID     base.NetworkID
}

func NewDigester(
	st *currencydigest.Database,
	root string,
	sourceReaders *isaac.BlockItemReaders,
	fromRemotes isaac.RemotesBlockItemReadFunc,
	networkID base.NetworkID,
	errChan chan error,
) *Digester {
	di := &Digester{
		Logging: logging.NewLogging(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "digester")
		}),
		database:      st,
		localfsRoot:   root,
		blockChan:     make(chan base.BlockMap, 100),
		errChan:       errChan,
		sourceReaders: sourceReaders,
		fromRemotes:   fromRemotes,
		networkID:     networkID,
	}

	di.ContextDaemon = util.NewContextDaemon(di.start)

	return di
}

func (di *Digester) start(ctx context.Context) error {
	e := util.StringError("start Digester")

	errch := func(err currencydigest.DigestError) {
		if di.errChan == nil {
			return
		}

		di.errChan <- err
	}

end:
	for {
		select {
		case <-ctx.Done():
			di.Log().Debug().Msg("stopped")

			break end
		case blk := <-di.blockChan:
			if m, _, _, _, _, _ := di.database.ManifestByHeight(blk.Manifest().Height()); m != nil {
				continue
			}

			err := util.Retry(ctx, func() (bool, error) {
				if err := di.digest(ctx, blk); err != nil {
					go errch(currencydigest.NewDigestError(err, blk.Manifest().Height()))

					if errors.Is(err, context.Canceled) {
						return false, e.Wrap(err)
					}

					return true, e.Wrap(err)
				}

				return false, nil
			}, 15, time.Second*1)
			if err != nil {
				di.Log().Error().Err(err).Int64("block", blk.Manifest().Height().Int64()).Msg("failed to digest block")
			} else {
				di.Log().Info().Int64("block", blk.Manifest().Height().Int64()).Msg("block digested")
			}

			go errch(currencydigest.NewDigestError(err, blk.Manifest().Height()))
		}
	}

	return nil
}

func (di *Digester) Digest(blocks []base.BlockMap) {
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Manifest().Height() < blocks[j].Manifest().Height()
	})

	for i := range blocks {
		blk := blocks[i]
		di.Log().Debug().Int64("block", blk.Manifest().Height().Int64()).Msg("start to digest block")

		di.blockChan <- blk
	}
}

func (di *Digester) digest(ctx context.Context, blk base.BlockMap) error {
	e := util.StringError("digest block")

	di.Lock()
	defer di.Unlock()

	var bm base.BlockMap

	switch i, found, err := isaac.BlockItemReadersDecode[base.BlockMap](di.sourceReaders.Item, blk.Manifest().Height(), base.BlockItemMap, nil); {
	case err != nil:
		return e.Wrap(err)
	case !found:
		return e.Wrap(util.ErrNotFound.Errorf("blockmap"))
	default:
		if err := i.IsValid(di.networkID); err != nil {
			return e.Wrap(err)
		}

		bm = i
	}

	pr, ops, sts, opsTree, _, _, err := isaacblock.LoadBlockItemsFromReader(bm, di.sourceReaders.Item, blk.Manifest().Height())
	if err != nil {
		return e.Wrap(err)
	}

	if err := DigestBlock(ctx, di.database, blk, ops, opsTree, sts, pr); err != nil {
		return e.Wrap(err)
	}

	return di.database.SetLastBlock(blk.Manifest().Height())
}

func DigestBlock(
	ctx context.Context,
	st *currencydigest.Database,
	blk base.BlockMap,
	ops []base.Operation,
	opstree fixedtree.Tree,
	sts []base.State,
	proposal base.ProposalSignFact,
) error {
	bs, err := NewBlockSession(st, blk, ops, opstree, sts, proposal)
	if err != nil {
		return err
	}
	defer func() {
		_ = bs.Close()
	}()

	if err := bs.Prepare(); err != nil {
		return err
	}

	return bs.Commit(ctx)
}
