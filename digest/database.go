package digest

import (
	"context"
	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	"strconv"

	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	nftstate "github.com/ProtoconNet/mitum-nft/v2/state"
	nfttypes "github.com/ProtoconNet/mitum-nft/v2/types"
	timestampservice "github.com/ProtoconNet/mitum-timestamp/state"
	timestamptypes "github.com/ProtoconNet/mitum-timestamp/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var maxLimit int64 = 50

var (
	defaultColNameAccount       = "digest_ac"
	defaultColNameBalance       = "digest_bl"
	defaultColNameCurrency      = "digest_cr"
	defaultColNameOperation     = "digest_op"
	defaultColNameBlock         = "digest_bm"
	defaultColNameNFTCollection = "digest_nftcollection"
	defaultColNameNFT           = "digest_nft"
	defaultColNameNFTOperator   = "digest_nftoperator"
	defaultColNameDIDIssuer     = "digest_did_issuer"
	defaultColNameDIDCredential = "digest_did_credential"
	defaultColNameHolderDID     = "digest_did_holder_did"
	defaultColNameTemplate      = "digest_did_template"
	defaultColNameTimeStamp     = "digest_ts"
)

func NFTCollection(st *currencydigest.Database, contract, col string) (*nfttypes.Design, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("collection", col)

	var design *nfttypes.Design
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameNFTCollection,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			design, err = nftstate.StateCollectionValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return design, nil
}

func NFT(st *currencydigest.Database, contract, col, idx string) (*nfttypes.NFT, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("collection", col)
	filter = filter.Add("nftid", idx)

	var nft *nfttypes.NFT
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameNFT,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			nft, err = nftstate.StateNFTValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return nft, nil
}

func NFTsByAddress(
	st *currencydigest.Database,
	address mitumbase.Address,
	reverse bool,
	offset string,
	limit int64,
	collectionid string,
	callback func(string /* nft id */, nfttypes.NFT) (bool, error),
) error {
	filter, err := buildNFTsFilterByAddress(address, offset, reverse, collectionid)
	if err != nil {
		return err
	}

	sr := 1
	if reverse {
		sr = -1
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("height", sr).D(),
	)

	switch {
	case limit <= 0: // no limit
	case limit > maxLimit:
		opt = opt.SetLimit(maxLimit)
	default:
		opt = opt.SetLimit(limit)
	}

	return st.DatabaseClient().Find(
		context.Background(),
		defaultColNameNFT,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.DatabaseEncoders())
			if err != nil {
				return false, err
			}
			nft, err := nftstate.StateNFTValue(st)
			if err != nil {
				return false, err
			}

			return callback(strconv.FormatUint(nft.ID(), 10), *nft)
		},
		opt,
	)
}

func NFTsByCollection(
	st *currencydigest.Database,
	contract,
	col string,
	reverse bool,
	offset string,
	limit int64,
	callback func(nft nfttypes.NFT, st mitumbase.State) (bool, error),
) error {
	filter, err := buildNFTsFilterByCollection(contract, col, offset, reverse)
	if err != nil {
		return err
	}

	sr := 1
	if reverse {
		sr = -1
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("height", sr).D(),
	)

	switch {
	case limit <= 0: // no limit
	case limit > maxLimit:
		opt = opt.SetLimit(maxLimit)
	default:
		opt = opt.SetLimit(limit)
	}

	return st.DatabaseClient().Find(
		context.Background(),
		defaultColNameNFT,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.DatabaseEncoders())
			if err != nil {
				return false, err
			}
			nft, err := nftstate.StateNFTValue(st)
			if err != nil {
				return false, err
			}
			return callback(*nft, st)
		},
		opt,
	)
}

func NFTOperators(
	st *currencydigest.Database,
	contract, col, account string,
) (*nfttypes.OperatorsBook, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("collection", col)
	filter = filter.Add("account", account)

	var operators *nfttypes.OperatorsBook
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameNFTOperator,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			operators, err = nftstate.StateOperatorsBookValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return operators, nil
}

func DIDService(st *currencydigest.Database, contract, svc string) (types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("service", svc)

	var design types.Design
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameDIDIssuer,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			design, err = state.StateDesignValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return types.Design{}, err
	}

	return design, nil
}

func Credential(st *currencydigest.Database, contract, svc, tid, id string) (types.Credential, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("service", svc)
	filter = filter.Add("template", tid)
	filter = filter.Add("credential_id", id)

	var credential types.Credential
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameDIDCredential,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			credential, err = state.StateCredentialValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return types.Credential{}, err
	}

	return credential, nil
}

func Template(st *currencydigest.Database, contract, svc, tid string) (types.Template, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("service", svc)
	filter = filter.Add("template", tid)

	var template types.Template
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameTemplate,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			template, err = state.StateTemplateValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return types.Template{}, err
	}

	return template, nil
}

func HolderDID(st *currencydigest.Database, contract, svc, address string) (string, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("service", svc)
	filter = filter.Add("holder", address)

	var did string
	var sta mitumbase.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameHolderDID,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			did, err = state.StateHolderDIDValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return "", err
	}

	return did, nil
}

func CredentialsByServiceAndTemplate(
	st *currencydigest.Database,
	contract,
	serviceID, templateID string,
	reverse bool,
	offset string,
	limit int64,
	callback func(types.Credential, mitumbase.State) (bool, error),
) error {
	filter, err := buildCredentialFilterByService(contract, serviceID, templateID, offset, reverse)
	if err != nil {
		return err
	}

	sr := 1
	if reverse {
		sr = -1
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("height", sr).D(),
	)

	switch {
	case limit <= 0: // no limit
	case limit > maxLimit:
		opt = opt.SetLimit(maxLimit)
	default:
		opt = opt.SetLimit(limit)
	}

	return st.DatabaseClient().Find(
		context.Background(),
		defaultColNameDIDCredential,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.DatabaseEncoders())
			if err != nil {
				return false, err
			}
			credential, err := state.StateCredentialValue(st)
			if err != nil {
				return false, err
			}
			return callback(credential, st)
		},
		opt,
	)
}

func buildCredentialFilterByService(contract, col, templateID string, offset string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
	filterSymbol := bson.D{{"service", bson.D{{"$in", []string{col}}}}}
	filterTemplate := bson.D{{"template", bson.D{{"$in", []string{templateID}}}}}
	filterA = append(filterA, filterSymbol)
	filterA = append(filterA, filterContract)
	filterA = append(filterA, filterTemplate)

	// if offset exist, apply offset
	if len(offset) > 0 {
		if !reverse {
			filterOffset := bson.D{
				{"credential_id", bson.D{{"$gt", offset}}},
			}
			filterA = append(filterA, filterOffset)
			// if reverse true, lesser then offset height
		} else {
			filterHeight := bson.D{
				{"credential_id", bson.D{{"$lt", offset}}},
			}
			filterA = append(filterA, filterHeight)
		}
	}

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	return filter, nil
}

func Timestamp(st *currencydigest.Database, contract, service string) (timestamptypes.Design, mitumbase.State, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("timestampservice", service)
	filter = filter.Add("isItem", false)
	q := filter.D()

	opt := options.FindOne().SetSort(
		util.NewBSONFilter("height", -1).D(),
	)
	var sta mitumbase.State
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameTimeStamp,
		q,
		func(res *mongo.SingleResult) error {
			i, err := currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			sta = i
			return nil
		},
		opt,
	); err != nil {
		return timestamptypes.Design{}, nil, err
	}

	if sta != nil {
		de, err := timestampservice.StateServiceDesignValue(sta)
		if err != nil {
			return timestamptypes.Design{}, nil, err
		}
		return de, sta, nil
	} else {
		return timestamptypes.Design{}, nil, errors.Errorf("state is nil")
	}
}

func TimestampItem(st *currencydigest.Database, contract, service, project string, idx uint64) (timestamptypes.TimeStampItem, mitumbase.State, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("timestampservice", service)
	filter = filter.Add("project", project)
	filter = filter.Add("timestampidx", idx)
	filter = filter.Add("isItem", true)
	q := filter.D()

	opt := options.FindOne().SetSort(
		util.NewBSONFilter("height", -1).D(),
	)
	var sta mitumbase.State
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameTimeStamp,
		q,
		func(res *mongo.SingleResult) error {
			i, err := currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			sta = i
			return nil
		},
		opt,
	); err != nil {
		return timestamptypes.TimeStampItem{}, nil, err
	}

	if sta != nil {
		it, err := timestampservice.StateTimeStampItemValue(sta)
		if err != nil {
			return timestamptypes.TimeStampItem{}, nil, err
		}
		return it, sta, nil
	} else {
		return timestamptypes.TimeStampItem{}, nil, errors.Errorf("state is nil")
	}
}
