package digest

import (
	mongodb "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bson "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

type TimeStampServiceDesignDoc struct {
	mongodb.BaseDoc
	st  base.State
	tsd types.Design
}

// NewTimeStampServiceDesignDoc gets the State of TimeStampServiceDesign
func NewTimeStampServiceDesignDoc(st base.State, enc encoder.Encoder) (TimeStampServiceDesignDoc, error) {
	tsd, err := state.StateServiceDesignValue(st)

	if err != nil {
		return TimeStampServiceDesignDoc{}, errors.Wrap(err, "TimeStampServiceDesignDoc needs ServiceDesign state")
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return TimeStampServiceDesignDoc{}, err
	}

	return TimeStampServiceDesignDoc{
		BaseDoc: b,
		st:      st,
		tsd:     tsd,
	}, nil
}

func (doc TimeStampServiceDesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.StateKeyTimeStampPrefix, 3)

	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()
	m["isItem"] = false

	return bson.Marshal(m)
}

type TimeStampItemDoc struct {
	mongodb.BaseDoc
	st     base.State
	tsItem types.TimeStampItem
}

func NewTimeStampItemDoc(st base.State, enc encoder.Encoder) (TimeStampItemDoc, error) {
	tsItem, err := state.StateTimeStampItemValue(st)
	if err != nil {
		return TimeStampItemDoc{}, errors.Wrap(err, "TimeStampServiceDesignDoc needs ServiceDesign state")
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return TimeStampItemDoc{}, err
	}

	return TimeStampItemDoc{
		BaseDoc: b,
		st:      st,
		tsItem:  tsItem,
	}, nil
}

func (doc TimeStampItemDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.StateKeyTimeStampPrefix, 5)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["project"] = doc.tsItem.ProjectID()
	m["timestampidx"] = doc.tsItem.TimestampID()
	m["height"] = doc.st.Height()
	m["isItem"] = true

	return bson.Marshal(m)
}
