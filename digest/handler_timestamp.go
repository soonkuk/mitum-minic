package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (hd *Handlers) handleTimeStamp(w http.ResponseWriter, r *http.Request) {
	cachekey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var contract string
	s, found := mux.Vars(r)["contract"]
	if !found {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusBadRequest)

		return
	}
	contract = s

	var service string
	s, found = mux.Vars(r)["service"]
	if !found {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusBadRequest)

		return
	}
	service = s

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleTimeStampInGroup(contract, service)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTimeStampInGroup(contract, service string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := Timestamp(hd.database, contract, service)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildTimeStamp(contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildTimeStamp(contract string, de types.Design, st base.State) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathTimeStampService, "contract", contract, "service", de.Service().String())
	if err != nil {
		return nil, err
	}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(de, currencydigest.NewHalLink(h, nil))

	h, err = hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handleTimeStampItem(w http.ResponseWriter, r *http.Request) {
	cachekey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var contract string
	s, found := mux.Vars(r)["contract"]
	if !found {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty contract address"), http.StatusBadRequest)

		return
	}
	contract = s

	var service string
	s, found = mux.Vars(r)["service"]
	if !found {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty service id"), http.StatusBadRequest)

		return
	}
	service = s

	var project string
	s, found = mux.Vars(r)["project"]
	if !found {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty project id"), http.StatusNotFound)

		return
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		currencydigest.HTTP2ProblemWithError(w, errors.Errorf("empty project id"), http.StatusBadRequest)

		return
	}
	project = s

	s, found = mux.Vars(r)["tid"]
	idx, err := parseIdxFromPath(s)
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleTimeStampItemInGroup(contract, service, project, idx)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTimeStampItemInGroup(contract, service, project string, idx uint64) ([]byte, error) {
	var it types.TimeStampItem
	var st base.State

	it, st, err := TimestampItem(hd.database, contract, service, project, idx)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildTimeStampItem(contract, service, it, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildTimeStampItem(contract, service string, it types.TimeStampItem, st base.State) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathTimeStampItem, "contract", contract, "service", service, "project", it.ProjectID(), "tid", strconv.FormatUint(it.TimestampID(), 10))
	if err != nil {
		return nil, err
	}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(it, currencydigest.NewHalLink(h, nil))

	h, err = hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}
