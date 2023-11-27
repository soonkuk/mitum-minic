package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mongodb "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bson "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	typesto "github.com/ProtoconNet/mitum-sto/types/sto"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type STODesignDoc struct {
	mongodb.BaseDoc
	st base.State
	de typesto.Design
}

func NewSTODesignDoc(st base.State, enc encoder.Encoder) (STODesignDoc, error) {
	de, err := ststo.StateDesignValue(st)
	if err != nil {
		return STODesignDoc{}, err
	}
	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STODesignDoc{}, err
	}

	return STODesignDoc{
		BaseDoc: b,
		st:      st,
		de:      de,
	}, nil
}

func (doc STODesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := ststo.ParseStateKey(doc.st.Key(), ststo.STOPrefix)
	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()
	m["design"] = doc.de

	return bson.Marshal(m)
}

type STOHolderPartitionsDoc struct {
	mongodb.BaseDoc
	st  base.State
	pts []typesto.Partition
}

func NewSTOHolderPartitionsDoc(st base.State, enc encoder.Encoder) (STOHolderPartitionsDoc, error) {
	pts, err := ststo.StateTokenHolderPartitionsValue(st)
	if err != nil {
		return STOHolderPartitionsDoc{}, err
	}
	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOHolderPartitionsDoc{}, err
	}

	return STOHolderPartitionsDoc{
		BaseDoc: b,
		st:      st,
		pts:     pts,
	}, nil
}

func (doc STOHolderPartitionsDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := ststo.ParseStateKey(doc.st.Key(), ststo.STOPrefix)
	m["contract"] = parsedKey[1]
	m["holder"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["partitions"] = doc.pts

	return bson.Marshal(m)
}

type STOHolderPartitionBalanceDoc struct {
	mongodb.BaseDoc
	st base.State
	am common.Big
}

func NewSTOHolderPartitionBalanceDoc(st base.State, enc encoder.Encoder) (STOHolderPartitionBalanceDoc, error) {
	am, err := ststo.StateTokenHolderPartitionBalanceValue(st)
	if err != nil {
		return STOHolderPartitionBalanceDoc{}, err
	}
	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOHolderPartitionBalanceDoc{}, err
	}

	return STOHolderPartitionBalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc STOHolderPartitionBalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := ststo.ParseStateKey(doc.st.Key(), ststo.STOPrefix)
	m["contract"] = parsedKey[1]
	m["holder"] = parsedKey[2]
	m["partition"] = parsedKey[3]
	m["height"] = doc.st.Height()
	m["balance"] = doc.am

	return bson.Marshal(m)
}

type STOHolderPartitionOperatorsDoc struct {
	mongodb.BaseDoc
	st   base.State
	oprs []base.Address
}

func NewSTOHolderPartitionOperatorsDoc(st base.State, enc encoder.Encoder) (STOHolderPartitionOperatorsDoc, error) {
	oprs, err := ststo.StateTokenHolderPartitionOperatorsValue(st)
	if err != nil {
		return STOHolderPartitionOperatorsDoc{}, err
	}
	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOHolderPartitionOperatorsDoc{}, err
	}

	return STOHolderPartitionOperatorsDoc{
		BaseDoc: b,
		st:      st,
		oprs:    oprs,
	}, nil
}

func (doc STOHolderPartitionOperatorsDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := ststo.ParseStateKey(doc.st.Key(), ststo.STOPrefix)
	m["contract"] = parsedKey[1]
	m["holder"] = parsedKey[2]
	m["partition"] = parsedKey[3]
	m["height"] = doc.st.Height()
	m["operators"] = doc.oprs

	return bson.Marshal(m)
}

type STOPartitionBalanceDoc struct {
	mongodb.BaseDoc
	st base.State
	am common.Big
}

func NewSTOPartitionBalanceDoc(st base.State, enc encoder.Encoder) (STOPartitionBalanceDoc, error) {
	am, err := ststo.StatePartitionBalanceValue(st)
	if err != nil {
		return STOPartitionBalanceDoc{}, err
	}
	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOPartitionBalanceDoc{}, err
	}

	return STOPartitionBalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc STOPartitionBalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := ststo.ParseStateKey(doc.st.Key(), ststo.STOPrefix)
	m["contract"] = parsedKey[1]
	m["partition"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["balance"] = doc.am

	return bson.Marshal(m)
}

type STOOperatorHoldersDoc struct {
	mongodb.BaseDoc
	st  base.State
	hds []base.Address
}

func NewSTOOperatorHoldersDoc(st base.State, enc encoder.Encoder) (STOOperatorHoldersDoc, error) {
	hds, err := ststo.StateOperatorTokenHoldersValue(st)
	if err != nil {
		return STOOperatorHoldersDoc{}, err
	}
	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return STOOperatorHoldersDoc{}, err
	}

	return STOOperatorHoldersDoc{
		BaseDoc: b,
		st:      st,
		hds:     hds,
	}, nil
}

func (doc STOOperatorHoldersDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := ststo.ParseStateKey(doc.st.Key(), ststo.STOPrefix)
	m["contract"] = parsedKey[1]
	m["operator"] = parsedKey[2]
	m["height"] = doc.st.Height()
	m["operators"] = doc.hds

	return bson.Marshal(m)
}
