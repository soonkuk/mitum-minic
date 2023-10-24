package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	timestampservice "github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Timestamp(st *currencydigest.Database, contract string) (types.Design, mitumbase.State, error) {
	filter := util.NewBSONFilter("contract", contract)
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
		return types.Design{}, nil, mitumutil.ErrNotFound.WithMessage(err, "timestamp service by contract %s", contract)
	}

	if sta != nil {
		de, err := timestampservice.StateServiceDesignValue(sta)
		if err != nil {
			return types.Design{}, nil, err
		}
		return de, sta, nil
	} else {
		return types.Design{}, nil, errors.Errorf("state is nil")
	}
}

func TimestampItem(st *currencydigest.Database, contract, project string, idx uint64) (types.TimeStampItem, mitumbase.State, error) {
	filter := util.NewBSONFilter("contract", contract)
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
		return types.TimeStampItem{}, nil, mitumutil.ErrNotFound.WithMessage(err, "timestamp item by contract %s, project %s, timestamp idx %s", contract, project, idx)
	}

	if sta != nil {
		it, err := timestampservice.StateTimeStampItemValue(sta)
		if err != nil {
			return types.TimeStampItem{}, nil, err
		}
		return it, sta, nil
	} else {
		return types.TimeStampItem{}, nil, errors.Errorf("state is nil")
	}
}
