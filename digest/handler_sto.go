package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	crcydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	typesto "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"net/http"
	"time"
)

func (hd *Handlers) handleSTOService(w http.ResponseWriter, r *http.Request) {
	cacheKey := crcydigest.CacheKeyPath(r)
	if err := crcydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleSTODesignInGroup(contract)
	}); err != nil {
		crcydigest.HTTP2HandleError(w, err)
	} else {
		crcydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			crcydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleSTODesignInGroup(contract string) (interface{}, error) {
	switch design, err := STOService(hd.database, contract); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(err, "sto service, contract %s", contract)
	case design == nil:
		return nil, util.ErrNotFound.Errorf("sto service, contract %s", contract)
	default:
		hal, err := hd.buildSTODesignHal(contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildSTODesignHal(contract string, design typesto.Design) (crcydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathSTOService, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := crcydigest.NewBaseHal(design, crcydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleSTOHolderPartitions(w http.ResponseWriter, r *http.Request) {
	cacheKey := crcydigest.CacheKeyPath(r)
	if err := crcydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	holder, err, status := parseRequest(w, r, "address")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleSTOHolderPartitionsInGroup(contract, holder)
	}); err != nil {
		crcydigest.HTTP2HandleError(w, err)
	} else {
		crcydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			crcydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleSTOHolderPartitionsInGroup(contract, holder string) (interface{}, error) {
	switch partitions, err := STOHolderPartitions(hd.database, contract, holder); {
	case err != nil:
		return nil, util.ErrNotFound.WithMessage(
			err,
			"partitions, contract %s, holder %s",
			contract,
			holder,
		)
	case partitions == nil:
		return nil, util.ErrNotFound.Errorf(
			"partitions, contract %s, holder %s",
			contract,
			holder,
		)
	default:
		hal, err := hd.buildSTOHolderPartitionsHal(contract, holder, partitions)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildSTOHolderPartitionsHal(
	contract, holder string, partitions []typesto.Partition,
) (crcydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathSTOHolderPartitions, "contract", contract, "address", holder)
	if err != nil {
		return nil, err
	}

	hal := crcydigest.NewBaseHal(partitions, crcydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleSTOHolderPartitionBalance(w http.ResponseWriter, r *http.Request) {
	cachekey := crcydigest.CacheKeyPath(r)
	if err := crcydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	holder, err, status := parseRequest(w, r, "address")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	partition, err, status := parseRequest(w, r, "partition")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleSTOHolderPartitionBalanceInGroup(contract, holder, partition)
	}); err != nil {
		crcydigest.HTTP2HandleError(w, err)
	} else {
		crcydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			crcydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleSTOHolderPartitionBalanceInGroup(
	contract, holder, partition string,
) (interface{}, error) {
	switch amount, err := STOHolderPartitionBalance(hd.database, contract, holder, partition); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildSTOHolderPartitionBalanceHal(contract, holder, partition, amount)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildSTOHolderPartitionBalanceHal(
	contract, holder, partition string, amount common.Big,
) (crcydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathSTOHolderPartitionBalance, "contract", contract, "address", holder, "partition", partition)
	if err != nil {
		return nil, err
	}

	hal := crcydigest.NewBaseHal(struct {
		Amount common.Big `json:"amount"`
	}{Amount: amount}, crcydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleSTOHolderPartitionOperators(w http.ResponseWriter, r *http.Request) {
	cachekey := crcydigest.CacheKeyPath(r)
	if err := crcydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	holder, err, status := parseRequest(w, r, "address")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	partition, err, status := parseRequest(w, r, "partition")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleSTOHolderPartitionOperatorsInGroup(contract, holder, partition)
	}); err != nil {
		crcydigest.HTTP2HandleError(w, err)
	} else {
		crcydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			crcydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleSTOHolderPartitionOperatorsInGroup(
	contract, holder, partition string,
) (interface{}, error) {
	switch operators, err := STOHolderPartitionOperators(hd.database, contract, holder, partition); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildSTOHolderPartitionOperatorsHal(contract, holder, partition, operators)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildSTOHolderPartitionOperatorsHal(
	contract, holder, partition string, operators []base.Address,
) (crcydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathSTOHolderPartitionOperators, "contract", contract, "address", holder, "partition", partition)
	if err != nil {
		return nil, err
	}

	hal := crcydigest.NewBaseHal(struct {
		Operators []base.Address `json:"operators"`
	}{Operators: operators}, crcydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleSTOPartitionBalance(w http.ResponseWriter, r *http.Request) {
	cachekey := crcydigest.CacheKeyPath(r)
	if err := crcydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	partition, err, status := parseRequest(w, r, "partition")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleSTOPartitionBalanceInGroup(contract, partition)
	}); err != nil {
		crcydigest.HTTP2HandleError(w, err)
	} else {
		crcydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			crcydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleSTOPartitionBalanceInGroup(
	contract, partition string,
) (interface{}, error) {
	switch amount, err := STOPartitionBalance(hd.database, contract, partition); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildSTOPartitionBalanceHal(contract, partition, amount)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildSTOPartitionBalanceHal(
	contract, partition string, amount common.Big,
) (crcydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathSTOPartitionBalance, "contract", contract, "partition", partition)
	if err != nil {
		return nil, err
	}

	hal := crcydigest.NewBaseHal(struct {
		Amount common.Big `json:"amount"`
	}{Amount: amount}, crcydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleSTOOperatorHolders(w http.ResponseWriter, r *http.Request) {
	cachekey := crcydigest.CacheKeyPath(r)
	if err := crcydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	operator, err, status := parseRequest(w, r, "address")
	if err != nil {
		crcydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleSTOOperatorHoldersInGroup(contract, operator)
	}); err != nil {
		crcydigest.HTTP2HandleError(w, err)
	} else {
		crcydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			crcydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleSTOOperatorHoldersInGroup(
	contract, operator string,
) (interface{}, error) {
	switch holders, err := STOOperatorHolders(hd.database, contract, operator); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildSTOOperatorHoldersHal(contract, operator, holders)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildSTOOperatorHoldersHal(
	contract, operator string, holders []base.Address,
) (crcydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathSTOOperatorHolders, "contract", contract, "address", operator)
	if err != nil {
		return nil, err
	}

	hal := crcydigest.NewBaseHal(struct {
		Holders []base.Address `json:"holders"`
	}{Holders: holders}, crcydigest.NewHalLink(h, nil))

	return hal, nil
}
