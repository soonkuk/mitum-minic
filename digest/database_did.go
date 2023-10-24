package digest

import (
	"context"

	"github.com/ProtoconNet/mitum-credential/state"
	"github.com/ProtoconNet/mitum-credential/types"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CredentialService(st *currencydigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta mitumbase.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameDIDCredentialService,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			de, err := state.StateDesignValue(sta)
			if err != nil {
				return err
			}
			design = &de

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return design, nil
}

func Credential(st *currencydigest.Database, contract, templateID, credentialID string) (*types.Credential, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("template", templateID)
	filter = filter.Add("credential_id", credentialID)

	var credential *types.Credential
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
			cre, err := state.StateCredentialValue(sta)
			if err != nil {
				return err
			}
			credential = &cre
			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return credential, nil
}

func Template(st *currencydigest.Database, contract, templateID string) (*types.Template, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("template", templateID)

	var template *types.Template
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
			te, err := state.StateTemplateValue(sta)
			if err != nil {
				return err
			}
			template = &te
			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return template, nil
}

func HolderDID(st *currencydigest.Database, contract, holder string) (string, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)

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
	templateID string,
	reverse bool,
	offset string,
	limit int64,
	callback func(types.Credential, mitumbase.State) (bool, error),
) error {
	filter, err := buildCredentialFilterByService(contract, templateID, offset, reverse)
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

func buildCredentialFilterByService(contract, templateID string, offset string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterContract := bson.D{{"contract", bson.D{{"$in", []string{contract}}}}}
	filterTemplate := bson.D{{"template", bson.D{{"$in", []string{templateID}}}}}
	filterA = append(filterA, filterContract)
	filterA = append(filterA, filterTemplate)

	// if offset exist, apply offset
	if len(offset) > 0 {
		if !reverse {
			filterOffset := bson.D{
				{"credential_id", bson.D{{"$gt", offset}}},
			}
			filterA = append(filterA, filterOffset)
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
