package cmds

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	credentialcmds "github.com/ProtoconNet/mitum-credential/cmds"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-minic/digest"
	nftcmds "github.com/ProtoconNet/mitum-nft/v2/cmds"
	timestampcmds "github.com/ProtoconNet/mitum-timestamp/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	isaacnetwork "github.com/ProtoconNet/mitum2/isaac/network"
	isaacstates "github.com/ProtoconNet/mitum2/isaac/states"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/network/quicmemberlist"
	"github.com/ProtoconNet/mitum2/network/quicstream"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/ProtoconNet/mitum2/util/ps"
	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type RunCommand struct { //nolint:govet //...
	launch.DesignFlag
	launch.DevFlags `embed:"" prefix:"dev."`
	Vault           string                `name:"vault" help:"privatekey path of vault"`
	Discovery       []launch.ConnInfoFlag `help:"member discovery" placeholder:"ConnInfo"`
	Hold            launch.HeightFlag     `help:"hold consensus states"`
	HTTPState       string                `name:"http-state" help:"runtime statistics thru https" placeholder:"bind address"`
	exitf           func(error)
	log             *zerolog.Logger
	holded          bool
}

func (cmd *RunCommand) Run(pctx context.Context) error {
	var log *logging.Logging
	if err := util.LoadFromContextOK(pctx, launch.LoggingContextKey, &log); err != nil {
		return err
	}

	log.Log().Debug().
		Interface("design", cmd.DesignFlag).
		Interface("vault", cmd.Vault).
		Interface("discovery", cmd.Discovery).
		Interface("hold", cmd.Hold).
		Interface("http_state", cmd.HTTPState).
		Interface("dev", cmd.DevFlags).
		Msg("flags")

	cmd.log = log.Log()

	if len(cmd.HTTPState) > 0 {
		if err := cmd.runHTTPState(cmd.HTTPState); err != nil {
			return errors.Wrap(err, "failed to run http state")
		}
	}

	//revive:disable:modifies-parameter
	pctx = context.WithValue(pctx, launch.DesignFlagContextKey, cmd.DesignFlag)
	pctx = context.WithValue(pctx, launch.DevFlagsContextKey, cmd.DevFlags)
	pctx = context.WithValue(pctx, launch.DiscoveryFlagContextKey, cmd.Discovery)
	pctx = context.WithValue(pctx, launch.VaultContextKey, cmd.Vault)
	//revive:enable:modifies-parameter

	pps := currencycmds.DefaultRunPS()

	_ = pps.AddOK(currencycmds.PNameMongoDBsDataBase, currencycmds.ProcessDatabase, nil, currencycmds.PNameDigestDesign, launch.PNameStorage).
		AddOK(currencycmds.PNameDigester, ProcessDigester, nil, currencycmds.PNameMongoDBsDataBase).
		AddOK(currencycmds.PNameDigest, currencycmds.ProcessDigestAPI, nil, currencycmds.PNameDigestDesign, currencycmds.PNameMongoDBsDataBase, launch.PNameMemberlist).
		AddOK(currencycmds.PNameDigestStart, currencycmds.ProcessStartDigestAPI, nil, currencycmds.PNameDigest).
		AddOK(currencycmds.PNameStartDigester, ProcessStartDigester, nil, currencycmds.PNameDigestStart)
	_ = pps.POK(launch.PNameStorage).PostAddOK(ps.Name("check-hold"), cmd.pCheckHold)
	_ = pps.POK(launch.PNameStates).
		PreAddOK(nftcmds.PNameOperationProcessorsMap, nftcmds.POperationProcessorsMap).
		PreAddOK(timestampcmds.PNameOperationProcessorsMap, timestampcmds.POperationProcessorsMap).
		PreAddOK(credentialcmds.PNameOperationProcessorsMap, credentialcmds.POperationProcessorsMap).
		PreAddOK(PNameOperationProcessorsMap, POperationProcessorsMap).
		PreAddOK(ps.Name("when-new-block-saved-in-consensus-state-func"), cmd.pWhenNewBlockSavedInConsensusStateFunc).
		PreAddOK(ps.Name("when-new-block-confirmed-func"), cmd.pWhenNewBlockConfirmed)
	_ = pps.POK(launch.PNameEncoder).
		PostAddOK(launch.PNameAddHinters, PAddHinters)
	_ = pps.POK(currencycmds.PNameDigest).
		PostAddOK(currencycmds.PNameDigestAPIHandlers, cmd.pDigestAPIHandlers)
	_ = pps.POK(currencycmds.PNameDigester).
		PostAddOK(currencycmds.PNameDigesterFollowUp, PDigesterFollowUp)

	_ = pps.SetLogging(log)

	log.Log().Debug().Interface("process", pps.Verbose()).Msg("process ready")

	pctx, err := pps.Run(pctx) //revive:disable-line:modifies-parameter
	defer func() {
		log.Log().Debug().Interface("process", pps.Verbose()).Msg("process will be closed")

		if _, err = pps.Close(pctx); err != nil {
			log.Log().Error().Err(err).Msg("failed to close")
		}
	}()

	if err != nil {
		return err
	}

	log.Log().Debug().
		Interface("discovery", cmd.Discovery).
		Interface("hold", cmd.Hold.Height()).
		Msg("node started")

	return cmd.run(pctx)
}

var errHoldStop = util.NewIDError("hold stop")

func (cmd *RunCommand) run(pctx context.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	exitch := make(chan error)

	cmd.exitf = func(err error) {
		exitch <- err
	}

	stopstates := func() {}

	if !cmd.holded {
		deferred, err := cmd.runStates(ctx, pctx)
		if err != nil {
			return err
		}

		stopstates = deferred
	}

	select {
	case <-ctx.Done(): // NOTE graceful stop
		return ctx.Err()
	case err := <-exitch:
		if errors.Is(err, errHoldStop) {
			stopstates()

			<-ctx.Done()

			return ctx.Err()
		}

		return err
	}
}

func (cmd *RunCommand) runStates(ctx, pctx context.Context) (func(), error) {
	var discoveries *util.Locked[[]quicstream.UDPConnInfo]
	var states *isaacstates.States

	if err := util.LoadFromContextOK(pctx,
		launch.DiscoveryContextKey, &discoveries,
		launch.StatesContextKey, &states,
	); err != nil {
		return nil, err
	}

	if dis := launch.GetDiscoveriesFromLocked(discoveries); len(dis) < 1 {
		cmd.log.Warn().Msg("empty discoveries; will wait to be joined by remote nodes")
	}

	go func() {
		cmd.exitf(<-states.Wait(ctx))
	}()

	return func() {
		if err := states.Hold(); err != nil && !errors.Is(err, util.ErrDaemonAlreadyStopped) {
			cmd.log.Error().Err(err).Msg("failed to stop states")

			return
		}

		cmd.log.Debug().Msg("states stopped")
	}, nil
}

func (cmd *RunCommand) pWhenNewBlockSavedInConsensusStateFunc(pctx context.Context) (context.Context, error) {
	var log *logging.Logging

	if err := util.LoadFromContextOK(pctx,
		launch.LoggingContextKey, &log,
	); err != nil {
		return pctx, err
	}

	f := func(bm base.BlockMap) {
		l := log.Log().With().
			Interface("blockmap", bm).
			Interface("height", bm.Manifest().Height()).
			Logger()

		if cmd.Hold.IsSet() && bm.Manifest().Height() == cmd.Hold.Height() {
			l.Debug().Msg("will be stopped by hold")

			cmd.exitf(errHoldStop.WithStack())

			return
		}
	}

	return context.WithValue(pctx, launch.WhenNewBlockSavedInConsensusStateFuncContextKey, f), nil
}

func (cmd *RunCommand) pWhenNewBlockConfirmed(pctx context.Context) (context.Context, error) {
	var log *logging.Logging
	var db isaac.Database
	var di *digest.Digester

	if err := util.LoadFromContextOK(pctx,
		launch.LoggingContextKey, &log,
		launch.CenterDatabaseContextKey, &db,
	); err != nil {
		return pctx, err
	}

	if err := util.LoadFromContext(pctx, currencycmds.ContextValueDigester, &di); err != nil {
		return pctx, err
	}

	var f func(height base.Height)
	if di != nil {
		g := cmd.whenBlockSaved(db, di)

		f = func(height base.Height) {
			g(pctx)
			l := log.Log().With().Interface("height", height).Logger()

			if cmd.Hold.IsSet() && height == cmd.Hold.Height() {
				l.Debug().Msg("will be stopped by hold")
				cmd.exitf(errHoldStop.WithStack())

				return
			}
		}
	} else {
		f = func(height base.Height) {
			l := log.Log().With().Interface("height", height).Logger()

			if cmd.Hold.IsSet() && height == cmd.Hold.Height() {
				l.Debug().Msg("will be stopped by hold")
				cmd.exitf(errHoldStop.WithStack())

				return
			}
		}
	}

	return context.WithValue(pctx,
		launch.WhenNewBlockConfirmedFuncContextKey, f,
	), nil
}

func (cmd *RunCommand) whenBlockSaved(
	db isaac.Database,
	di *digest.Digester,
) ps.Func {
	return func(ctx context.Context) (context.Context, error) {
		switch m, found, err := db.LastBlockMap(); {
		case err != nil:
			return ctx, err
		case !found:
			return ctx, errors.Errorf("last BlockMap not found")
		default:
			if di != nil {
				go func() {
					di.Digest([]base.BlockMap{m})
				}()
			}
		}
		return ctx, nil
	}
}

func (cmd *RunCommand) pCheckHold(pctx context.Context) (context.Context, error) {
	var db isaac.Database
	if err := util.LoadFromContextOK(pctx, launch.CenterDatabaseContextKey, &db); err != nil {
		return pctx, err
	}

	switch {
	case !cmd.Hold.IsSet():
	case cmd.Hold.Height() < base.GenesisHeight:
		cmd.holded = true
	default:
		switch m, found, err := db.LastBlockMap(); {
		case err != nil:
			return pctx, err
		case !found:
		case cmd.Hold.Height() <= m.Manifest().Height():
			cmd.holded = true
		}
	}

	return pctx, nil
}

func (cmd *RunCommand) runHTTPState(bind string) error {
	addr, err := net.ResolveTCPAddr("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "failed to parse --http-state")
	}

	m := http.NewServeMux()
	if err := statsviz.Register(m); err != nil {
		return errors.Wrap(err, "failed to register statsviz for http-state")
	}

	cmd.log.Debug().Stringer("bind", addr).Msg("statsviz started")

	go func() {
		_ = http.ListenAndServe(addr.String(), m)
	}()

	return nil
}

func (cmd *RunCommand) pDigestAPIHandlers(ctx context.Context) (context.Context, error) {
	var isaacparams *isaac.Params
	var local base.LocalNode

	if err := util.LoadFromContextOK(ctx,
		launch.LocalContextKey, &local,
		launch.ISAACParamsContextKey, &isaacparams,
	); err != nil {
		return nil, err
	}

	var design currencycmds.DigestDesign
	if err := util.LoadFromContext(ctx, currencycmds.ContextValueDigestDesign, &design); err != nil {
		if errors.Is(err, util.ErrNotFound) {
			return ctx, nil
		}

		return nil, err
	}

	if (design == currencycmds.DigestDesign{}) {
		return ctx, nil
	}

	cache, err := cmd.loadCache(ctx, design)
	if err != nil {
		return ctx, err
	}

	var dnt *currencydigest.HTTP2Server
	if err := util.LoadFromContext(ctx, currencycmds.ContextValueDigestNetwork, &dnt); err != nil {
		return ctx, err
	}
	router := dnt.Router()

	defaultHandlers, err := cmd.setDigestDefaultHandlers(ctx, isaacparams, cache, router)
	if err != nil {
		return ctx, err
	}

	if err := defaultHandlers.Initialize(); err != nil {
		return ctx, err
	}

	handlers, err := cmd.setDigestHandlers(ctx, isaacparams, cache, router)
	if err != nil {
		return ctx, err
	}

	if err := handlers.Initialize(); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (cmd *RunCommand) loadCache(_ context.Context, design currencycmds.DigestDesign) (currencydigest.Cache, error) {
	c, err := currencydigest.NewCacheFromURI(design.Cache().String())
	if err != nil {
		cmd.log.Err(err).Str("cache", design.Cache().String()).Msg("failed to connect cache server")
		cmd.log.Warn().Msg("instead of remote cache server, internal mem cache can be available, `memory://`")

		return nil, err
	}
	return c, nil
}

func (cmd *RunCommand) setDigestDefaultHandlers(
	ctx context.Context,
	params *isaac.Params,
	cache currencydigest.Cache,
	router *mux.Router,
) (*currencydigest.Handlers, error) {
	var st *currencydigest.Database
	if err := util.LoadFromContext(ctx, currencycmds.ContextValueDigestDatabase, &st); err != nil {
		return nil, err
	}

	handlers := currencydigest.NewHandlers(ctx, params.NetworkID(), encs, enc, st, cache, router)

	h, err := cmd.setDigestNetworkClient(ctx, params, handlers)
	if err != nil {
		return nil, err
	}
	handlers = h

	return handlers, nil
}

func (cmd *RunCommand) setDigestHandlers(
	ctx context.Context,
	params *isaac.Params,
	cache currencydigest.Cache,
	router *mux.Router,
) (*digest.Handlers, error) {
	var st *currencydigest.Database
	if err := util.LoadFromContext(ctx, currencycmds.ContextValueDigestDatabase, &st); err != nil {
		return nil, err
	}

	handlers := digest.NewHandlers(ctx, params.NetworkID(), encs, enc, st, cache, router)

	return handlers, nil
}

func (cmd *RunCommand) setDigestNetworkClient(
	ctx context.Context,
	params *isaac.Params,
	handlers *currencydigest.Handlers,
) (*currencydigest.Handlers, error) {
	var memberList *quicmemberlist.Memberlist
	if err := util.LoadFromContextOK(ctx, launch.MemberlistContextKey, &memberList); err != nil {
		return nil, err
	}

	client := launch.NewNetworkClient( //nolint:gomnd //...
		encs, enc,
		(*params).NetworkID(),
	)

	handlers = handlers.SetNetworkClientFunc(
		func() (*isaacnetwork.QuicstreamClient, *quicmemberlist.Memberlist, error) { // nolint:contextcheck
			return client, memberList, nil
		},
	)

	cmd.log.Debug().Msg("send handler attached")

	return handlers, nil
}
