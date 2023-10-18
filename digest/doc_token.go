package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	mongodbstorage "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type TokenDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	de types.Design
}

func NewTokenDoc(st base.State, enc encoder.Encoder) (TokenDoc, error) {
	de, err := state.StateDesignValue(st)
	if err != nil {
		return TokenDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return TokenDoc{}, err
	}

	return TokenDoc{
		BaseDoc: b,
		st:      st,
		de:      *de,
	}, nil
}

func (doc TokenDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	stateKeys, err := state.ParseStateKey(doc.st.Key(), state.TokenPrefix)
	if err != nil {
		return nil, err
	}
	m["contract"] = stateKeys[1]
	m["height"] = doc.st.Height()
	m["design"] = doc.de

	return bsonenc.Marshal(m)
}

type BalanceDoc struct {
	mongodbstorage.BaseDoc
	st     base.State
	amount common.Big
}

func NewBalanceDoc(st base.State, enc encoder.Encoder) (*BalanceDoc, error) {
	balance, err := state.StateTokenBalanceValue(st)
	if err != nil {
		return nil, err
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &BalanceDoc{
		BaseDoc: b,
		st:      st,
		amount:  balance,
	}, nil
}

func (doc BalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	stateKeys, err := state.ParseStateKey(doc.st.Key(), state.TokenPrefix)
	if err != nil {
		return nil, err
	}
	m["contract"] = stateKeys[1]
	m["address"] = stateKeys[2]
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
