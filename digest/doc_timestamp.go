package digest

import (
	mongodbstorage "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	timestampservice "github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

type TimeStampServiceDesignDoc struct {
	mongodbstorage.BaseDoc
	st  base.State
	tsd types.Design
}

// NewTimeStampServiceDesignDoc gets the State of TimeStampServiceDesign
func NewTimeStampServiceDesignDoc(st base.State, enc encoder.Encoder) (TimeStampServiceDesignDoc, error) {
	tsd, err := timestampservice.StateServiceDesignValue(st)

	if err != nil {
		return TimeStampServiceDesignDoc{}, errors.Wrap(err, "TimeStampServiceDesignDoc needs ServiceDesign state")
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
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

	parsedKey, err := timestampservice.ParseStateKey(doc.st.Key())

	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()
	m["isItem"] = false

	return bsonenc.Marshal(m)
}

type TimeStampItemDoc struct {
	mongodbstorage.BaseDoc
	st     base.State
	tsItem types.TimeStampItem
}

func NewTimeStampItemDoc(st base.State, enc encoder.Encoder) (TimeStampItemDoc, error) {
	tsItem, err := timestampservice.StateTimeStampItemValue(st)
	if err != nil {
		return TimeStampItemDoc{}, errors.Wrap(err, "TimeStampServiceDesignDoc needs ServiceDesign state")
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
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

	parsedKey, err := timestampservice.ParseStateKey(doc.st.Key())
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["project"] = doc.tsItem.ProjectID()
	m["timestampidx"] = doc.tsItem.TimestampID()
	m["height"] = doc.st.Height()
	m["isItem"] = true

	return bsonenc.Marshal(m)
}
