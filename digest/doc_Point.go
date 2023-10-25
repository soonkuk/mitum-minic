package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mongodbstorage "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-point/state"
	"github.com/ProtoconNet/mitum-point/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type PointDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	de types.Design
}

func NewPointDoc(st base.State, enc encoder.Encoder) (PointDoc, error) {
	de, err := state.StateDesignValue(st)
	if err != nil {
		return PointDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return PointDoc{}, err
	}

	return PointDoc{
		BaseDoc: b,
		st:      st,
		de:      *de,
	}, nil
}

func (doc PointDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	stateKeys, err := state.ParseStateKey(doc.st.Key(), state.PointPrefix)
	if err != nil {
		return nil, err
	}
	m["contract"] = stateKeys[1]
	m["height"] = doc.st.Height()
	m["design"] = doc.de

	return bsonenc.Marshal(m)
}

type PointBalanceDoc struct {
	mongodbstorage.BaseDoc
	st     base.State
	amount common.Big
}

func NewPointBalanceDoc(st base.State, enc encoder.Encoder) (*PointBalanceDoc, error) {
	balance, err := state.StatePointBalanceValue(st)
	if err != nil {
		return nil, err
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &PointBalanceDoc{
		BaseDoc: b,
		st:      st,
		amount:  balance,
	}, nil
}

func (doc PointBalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	stateKeys, err := state.ParseStateKey(doc.st.Key(), state.PointPrefix)
	if err != nil {
		return nil, err
	}
	m["contract"] = stateKeys[1]
	m["address"] = stateKeys[2]
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
