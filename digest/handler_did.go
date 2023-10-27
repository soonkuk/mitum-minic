package digest

import (
	"github.com/ProtoconNet/mitum-credential/types"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/ProtoconNet/mitum2/base"
)

func (hd *Handlers) handleCredentialService(w http.ResponseWriter, r *http.Request) {
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
		return hd.handleCredentialServiceInGroup(contract)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleCredentialServiceInGroup(contract string) (interface{}, error) {
	switch design, err := CredentialService(hd.database, contract); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "credential service, contract %s", contract)
	case design == nil:
		return nil, mitumutil.ErrNotFound.Errorf("credential service, contract %s", contract)
	default:
		hal, err := hd.buildCredentialServiceHal(contract, *design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildCredentialServiceHal(contract string, design types.Design) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDIDService, "contract", contract)
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
	switch credential, isActive, err := Credential(hd.database, contract, templateID, credentialID); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "credential by contract %s, template %s, id %s", contract, templateID, credentialID)
	case credential == nil:
		return nil, mitumutil.ErrNotFound.Errorf("credential by contract %s, template %s, id %s", contract, templateID, credentialID)
	default:
		hal, err := hd.buildCredentialHal(contract, *credential, isActive)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildCredentialHal(
	contract string,
	credential types.Credential,
	isActive bool,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDIDCredential,
		"contract", contract,
		"templateid", credential.TemplateID(),
		"credentialid", credential.ID(),
	)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(
		struct {
			Credential types.Credential `json:"credential"`
			IsActive   bool             `json:"is_active"`
		}{Credential: credential, IsActive: isActive},
		currencydigest.NewHalLink(h, nil),
	)

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
		func(credential types.Credential, isActive bool, st base.State) (bool, error) {
			hal, err := hd.buildCredentialHal(contract, credential, isActive)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, mitumutil.ErrNotFound.WithMessage(err, "credentials by contract %s, template %s", contract, templateID)
	} else if len(vas) < 1 {
		return nil, false, mitumutil.ErrNotFound.Errorf("credentials by contract %s, template %s", contract, templateID)
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

	h, err := hd.combineURL(HandlerPathDIDService, "contract", contract)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("service", currencydigest.NewHalLink(h, nil))

	var nextOffset string

	if len(vas) > 0 {
		va, ok := vas[len(vas)-1].Interface().(struct {
			Credential types.Credential `json:"credential"`
			IsActive   bool             `json:"is_active"`
		})
		if !ok {
			return nil, errors.Errorf("failed to build credentials hal")
		}
		nextOffset = va.Credential.ID()

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

func (hd *Handlers) handleHolderCredential(w http.ResponseWriter, r *http.Request) {
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
		return hd.handleHolderCredentialsInGroup(contract, holder)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleHolderCredentialsInGroup(contract, holder string) (interface{}, error) {
	var did string
	switch d, err := HolderDID(hd.database, contract, holder); {
	case err != nil:
		return nil, mitumutil.ErrNotFound.WithMessage(err, "DID by contract %s, holder %s", contract, holder)
	case d == "":
		return nil, mitumutil.ErrNotFound.Errorf("DID by contract %s, holder %s", contract, holder)
	default:
		did = d
	}

	var vas []currencydigest.Hal
	if err := CredentialsByServiceHolder(
		hd.database, contract, holder,
		func(credential types.Credential, isActive bool, st base.State) (bool, error) {
			hal, err := hd.buildCredentialHal(contract, credential, isActive)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, mitumutil.ErrNotFound.WithMessage(err, "credentials by contract %s, holder %s", contract, holder)
	} else if len(vas) < 1 {
		return nil, mitumutil.ErrNotFound.Errorf("credentials by contract %s, holder %s", contract, holder)
	}
	hal, err := hd.buildHolderDIDCredentialsHal(contract, holder, did, vas)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(hal)
}

func (hd *Handlers) buildHolderDIDCredentialsHal(
	contract, holder, did string,
	vas []currencydigest.Hal,
) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDIDHolder, "contract", contract, "holder", holder)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(
		struct {
			DID         string               `json:"did"`
			Credentials []currencydigest.Hal `json:"credentials"`
		}{
			DID:         did,
			Credentials: vas,
		}, currencydigest.NewHalLink(h, nil))

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
		return nil, mitumutil.ErrNotFound.WithMessage(err, "template by contract %s, template %s", contract, templateID)
	case template == nil:
		return nil, mitumutil.ErrNotFound.Errorf("template by contract %s, template %s", contract, templateID)
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
