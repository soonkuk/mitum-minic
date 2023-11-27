package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	crcydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	typesto "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func STOService(
	st *crcydigest.Database,
	contract string,
) (*typesto.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design typesto.Design
	var sta base.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameSTO,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = crcydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			design, err = ststo.StateDesignValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return &design, nil
}

func STOHolderPartitions(
	st *crcydigest.Database,
	contract,
	holder string,
) ([]typesto.Partition, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)

	var partitions []typesto.Partition
	var sta base.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameSTOHolderPartitions,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = crcydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			partitions, err = ststo.StateTokenHolderPartitionsValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return partitions, nil
}

func STOHolderPartitionBalance(
	st *crcydigest.Database,
	contract,
	holder,
	partition string,
) (common.Big, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)
	filter = filter.Add("partition", partition)

	var amount common.Big
	var sta base.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameSTOHolderPartitionBalance,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = crcydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			amount, err = ststo.StateTokenHolderPartitionBalanceValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return common.NilBig, mitumutil.ErrNotFound.Errorf(
			"sto holder partition balance by contract %s, account %s",
			contract,
			holder,
		)
	}

	return amount, nil
}

func STOHolderPartitionOperators(
	st *crcydigest.Database,
	contract,
	holder,
	partition string,
) ([]base.Address, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("holder", holder)
	filter = filter.Add("partition", partition)

	var operators []base.Address
	var sta base.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameSTOHolderPartitionOperators,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = crcydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			operators, err = ststo.StateTokenHolderPartitionOperatorsValue(sta)
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

func STOPartitionBalance(
	st *crcydigest.Database,
	contract,
	partition string,
) (common.Big, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("partition", partition)

	var amount common.Big
	var sta base.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameSTOPartitionBalance,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = crcydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			amount, err = ststo.StatePartitionBalanceValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return common.NilBig, mitumutil.ErrNotFound.Errorf(
			"sto partition balance by contract %s, account %s",
			contract,
			partition,
		)
	}

	return amount, nil
}

func STOOperatorHolders(
	st *crcydigest.Database,
	contract,
	operator string,
) ([]base.Address, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("operator", operator)

	var holders []base.Address
	var sta base.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameSTOOperatorHolders,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = crcydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			holders, err = ststo.StateOperatorTokenHoldersValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return holders, nil
}
