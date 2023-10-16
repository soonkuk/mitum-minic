package digest

import (
	"context"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NFTCollection(st *currencydigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
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

			design, err = state.StateCollectionValue(sta)
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

func NFT(st *currencydigest.Database, contract, idx string) (*types.NFT, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("nftid", idx)

	var nft *types.NFT
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
			nft, err = state.StateNFTValue(sta)
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

func NFTsByCollection(
	st *currencydigest.Database,
	contract string,
	reverse bool,
	offset string,
	limit int64,
	callback func(nft types.NFT, st mitumbase.State) (bool, error),
) error {
	filter, err := buildNFTsFilterByContract(contract, offset, reverse)
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
			nft, err := state.StateNFTValue(st)
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
	contract, account string,
) (*types.OperatorsBook, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("address", account)

	var operators *types.OperatorsBook
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

			operators, err = state.StateOperatorsBookValue(sta)
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
