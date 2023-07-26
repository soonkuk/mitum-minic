package digest

import (
	"github.com/ProtoconNet/mitum-nft/v2/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

func (bs *BlockSession) prepareNFTs() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var nftCollectionModels []mongo.WriteModel
	var nftOperatorModels []mongo.WriteModel
	var nftBoxModels []mongo.WriteModel
	var nftModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]
		stateKey, err := state.ParseNFTStateKey(st.Key())
		if err != nil {
			continue
		}
		switch stateKey {
		case state.CollectionKey:
			j, err := bs.handleNFTCollectionState(st)
			if err != nil {
				return err
			}
			nftCollectionModels = append(nftCollectionModels, j...)
		case state.OperatorsKey:
			j, err := bs.handleNFTOperatorsState(st)
			if err != nil {
				return err
			}
			nftOperatorModels = append(nftOperatorModels, j...)
		case state.NFTBoxKey:
			j, err := bs.handleNFTBoxState(st)
			if err != nil {
				return err
			}
			nftBoxModels = append(nftBoxModels, j...)
		case state.NFTKey:
			j, nft, err := bs.handleNFTState(st)
			if err != nil {
				return err
			}
			nftModels = append(nftModels, j...)
			bs.nftMap[strconv.FormatUint(nft, 10)] = struct{}{}
		default:
			continue
		}
	}

	bs.nftCollectionModels = nftCollectionModels
	bs.nftOperatorModels = nftOperatorModels
	bs.nftBoxModels = nftBoxModels
	bs.nftModels = nftModels

	return nil
}

func (bs *BlockSession) handleNFTOperatorsState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftCollectionDoc, err := NewNFTOperatorDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftCollectionDoc),
		}, nil
	}
}

func (bs *BlockSession) handleNFTState(st mitumbase.State) ([]mongo.WriteModel, uint64, error) {
	if nftDoc, err := NewNFTDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, 0, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftDoc),
		}, nftDoc.nft.ID(), nil
	}
}

func (bs *BlockSession) handleNFTBoxState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftBoxDoc, err := NewNFTBoxDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftBoxDoc),
		}, nil
	}
}

func (bs *BlockSession) handleNFTLastIndexState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if nftLastIndexDoc, err := NewNFTLastIndexDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(nftLastIndexDoc),
		}, nil
	}
}
