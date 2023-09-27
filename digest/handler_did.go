package digest

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"net/http"
	"time"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

func (hd *Handlers) handleDIDIssuer(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDIDIssuerInGroup(contract)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDIDIssuerInGroup(contract string) (interface{}, error) {
	switch design, err := DIDService(hd.database, contract); {
	case err != nil:
		return nil, err
	case design == nil:
		return nil, mitumutil.ErrNotFound.Errorf("issuer design, %v in handleDIDIssuer contract", contract)
	default:
		hal, err := hd.buildDIDServiceHal(contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildDIDServiceHal(contract string, design types.Design) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDIDIssuer, "contract", contract)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(design, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleCredential(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	templateID, err, status := parseRequest(w, r, "templateid")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	credentialID, err, status := parseRequest(w, r, "credentialid")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleCredentialInGroup(contract, templateID, credentialID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleCredentialInGroup(contract, templateID, credentialID string) (interface{}, error) {
	switch credential, err := Credential(hd.database, contract, templateID, credentialID); {
	case err != nil:
		return nil, err
	case credential == nil:
		return nil, mitumutil.ErrNotFound.Errorf("credential, %v in handleCredential", credentialID)
	default:
		hal, err := hd.buildCredentialHal(contract, templateID, *credential)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildCredentialHal(
	contract, templateID string,
	credential types.Credential,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDIDCredential,
		"contract", contract,
		"templateid", templateID,
		"credentialid", credential.ID(),
	)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(credential, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleCredentials(w http.ResponseWriter, r *http.Request) {
	limit := currencydigest.ParseLimitQuery(r.URL.Query().Get("limit"))
	offset := currencydigest.ParseStringQuery(r.URL.Query().Get("offset"))
	reverse := currencydigest.ParseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := currencydigest.CacheKey(
		r.URL.Path, currencydigest.StringOffsetQuery(offset),
		currencydigest.StringBoolQuery("reverse", reverse),
	)

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	templateID, err, status := parseRequest(w, r, "templateid")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleCredentialsInGroup(contract, templateID, offset, reverse, limit)

		return []interface{}{i, filled}, err
	})

	if err != nil {
		hd.Log().Err(err).Str("Issuer", contract).Msg("failed to get credentials")
		currencydigest.HTTP2HandleError(w, err)

		return
	}

	var b []byte
	var filled bool
	{
		l := v.([]interface{})
		b = l[0].([]byte)
		filled = l[1].(bool)
	}

	currencydigest.HTTP2WriteHalBytes(hd.encoder, w, b, http.StatusOK)

	if !shared {
		expire := hd.expireNotFilled
		if len(offset) > 0 && filled {
			expire = time.Minute
		}

		currencydigest.HTTP2WriteCache(w, cachekey, expire)
	}
}

func (hd *Handlers) handleCredentialsInGroup(
	contract, templateID string,
	offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("service-credentials")
	} else {
		limit = l
	}

	var vas []currencydigest.Hal
	if err := CredentialsByServiceAndTemplate(
		hd.database, contract, templateID, reverse, offset, limit,
		func(credential types.Credential, st base.State) (bool, error) {
			hal, err := hd.buildCredentialHal(contract, templateID, credential)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, err
	} else if len(vas) < 1 {
		return nil, false, errors.Errorf("credentials not found")
	}

	i, err := hd.buildCredentialsHal(contract, templateID, vas, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	b, err := hd.encoder.Marshal(i)
	return b, int64(len(vas)) == limit, err
}

func (hd *Handlers) buildCredentialsHal(
	contract, templateID string,
	vas []currencydigest.Hal,
	offset string,
	reverse bool,
) (currencydigest.Hal, error) {
	baseSelf, err := hd.combineURL(
		HandlerPathDIDCredentials,
		"contract", contract,
		"templateid", templateID,
	)
	if err != nil {
		return nil, err
	}

	self := baseSelf
	if len(offset) > 0 {
		self = currencydigest.AddQueryValue(baseSelf, currencydigest.StringOffsetQuery(offset))
	}
	if reverse {
		self = currencydigest.AddQueryValue(baseSelf, currencydigest.StringBoolQuery("reverse", reverse))
	}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(vas, currencydigest.NewHalLink(self, nil))

	h, err := hd.combineURL(HandlerPathDIDIssuer, "contract", contract)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("service", currencydigest.NewHalLink(h, nil))

	var nextOffset string

	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(types.Credential)
		nextOffset = va.ID()
	}

	if len(nextOffset) > 0 {
		next := baseSelf
		next = currencydigest.AddQueryValue(next, currencydigest.StringOffsetQuery(nextOffset))

		if reverse {
			next = currencydigest.AddQueryValue(next, currencydigest.StringBoolQuery("reverse", reverse))
		}

		hal = hal.AddLink("next", currencydigest.NewHalLink(next, nil))
	}

	hal = hal.AddLink("reverse", currencydigest.NewHalLink(currencydigest.AddQueryValue(baseSelf, currencydigest.StringBoolQuery("reverse", !reverse)), nil))

	return hal, nil
}

func (hd *Handlers) handleHolderDID(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	holder, err, status := parseRequest(w, r, "holder")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleHolderDIDInGroup(contract, holder)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleHolderDIDInGroup(contract, holder string) (interface{}, error) {
	switch did, err := HolderDID(hd.database, contract, holder); {
	case err != nil:
		return nil, err
	case did == "":
		return nil, mitumutil.ErrNotFound.Errorf("DID for holder, %v in handleHolderDID", holder)
	default:
		hal, err := hd.buildHolderDIDHal(contract, holder, did)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildHolderDIDHal(
	contract, holder, did string,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDIDHolder, "contract", contract, "holder", holder)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(did, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleTemplate(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	templateID, err, status := parseRequest(w, r, "templateid")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleTemplateInGroup(contract, templateID)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTemplateInGroup(contract, templateID string) (interface{}, error) {
	switch template, err := Template(hd.database, contract, templateID); {
	case err != nil:
		return nil, err
	case template == nil:
		return nil, mitumutil.ErrNotFound.Errorf("template, %v in handleTemplate", templateID)
	default:
		hal, err := hd.buildTemplateHal(contract, templateID, *template)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildTemplateHal(
	contract, templateID string,
	template types.Template,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDIDTemplate,
		"contract", contract,
		"templateid", templateID,
	)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(template, currencydigest.NewHalLink(h, nil))

	return hal, nil
}
