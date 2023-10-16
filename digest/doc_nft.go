package digest

import (
	mongodbstorage "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"strconv"
)

type NFTCollectionDoc struct {
	mongodbstorage.BaseDoc
	st base.State
	de types.Design
}

func NewNFTCollectionDoc(st base.State, enc encoder.Encoder) (NFTCollectionDoc, error) {
	de, err := state.StateCollectionValue(st)
	if err != nil {
		return NFTCollectionDoc{}, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return NFTCollectionDoc{}, err
	}

	return NFTCollectionDoc{
		BaseDoc: b,
		st:      st,
		de:      *de,
	}, nil
}

func (doc NFTCollectionDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["contract"] = doc.de.Parent()
	m["height"] = doc.st.Height()
	m["design"] = doc.de

	return bsonenc.Marshal(m)
}

type NFTDoc struct {
	mongodbstorage.BaseDoc
	st        base.State
	nft       types.NFT
	addresses []base.Address
	owner     string
}

func NewNFTDoc(st base.State, enc encoder.Encoder) (*NFTDoc, error) {
	nft, err := state.StateNFTValue(st)
	if err != nil {
		return nil, err
	}
	var addresses = make([]string, len(nft.Creators().Addresses())+1)
	addresses[0] = nft.Owner().String()
	for i := range nft.Creators().Addresses() {
		addresses[i+1] = nft.Creators().Addresses()[i].String()
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &NFTDoc{
		BaseDoc:   b,
		st:        st,
		nft:       *nft,
		addresses: nft.Addresses(),
		owner:     nft.Owner().String(),
	}, nil
}

func (doc NFTDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.NFTPrefix)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["nftid"] = strconv.FormatUint(doc.nft.ID(), 10)
	m["owner"] = doc.nft.Owner()
	m["addresses"] = doc.addresses
	m["istoken"] = true
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type NFTOperatorDoc struct {
	mongodbstorage.BaseDoc
	st        base.State
	operators types.OperatorsBook
}

func NewNFTOperatorDoc(st base.State, enc encoder.Encoder) (*NFTOperatorDoc, error) {
	operators, err := state.StateOperatorsBookValue(st)
	if err != nil {
		return nil, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &NFTOperatorDoc{
		BaseDoc:   b,
		st:        st,
		operators: *operators,
	}, nil
}

func (doc NFTOperatorDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.NFTPrefix)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["address"] = parsedKey[2]
	m["operators"] = doc.operators
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type NFTBoxDoc struct {
	mongodbstorage.BaseDoc
	st     base.State
	nftbox types.NFTBox
}

func NewNFTBoxDoc(st base.State, enc encoder.Encoder) (*NFTBoxDoc, error) {
	nftbox, err := state.StateNFTBoxValue(st)
	if err != nil {
		return nil, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &NFTBoxDoc{
		BaseDoc: b,
		st:      st,
		nftbox:  nftbox,
	}, nil
}

func (doc NFTBoxDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.NFTPrefix)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["nfts"] = doc.nftbox.NFTs()
	m["istoken"] = false
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type NFTLastIndexDoc struct {
	mongodbstorage.BaseDoc
	st    base.State
	nftID uint64
}

func NewNFTLastIndexDoc(st base.State, enc encoder.Encoder) (*NFTLastIndexDoc, error) {
	nftID, err := state.StateLastNFTIndexValue(st)
	if err != nil {
		return nil, err
	}
	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &NFTLastIndexDoc{
		BaseDoc: b,
		st:      st,
		nftID:   nftID,
	}, nil
}

func (doc NFTLastIndexDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	parsedKey, err := state.ParseStateKey(doc.st.Key(), state.NFTPrefix)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["id"] = doc.nftID
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
