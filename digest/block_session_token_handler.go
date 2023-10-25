package digest

import (
	"github.com/ProtoconNet/mitum-token/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareToken() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var TokenModels []mongo.WriteModel
	var TokenBalanceModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]

		switch {
		case state.IsStateDesignKey(st.Key()):
			j, err := bs.handleTokenState(st)
			if err != nil {
				return err
			}
			TokenModels = append(TokenModels, j...)
		case state.IsStateTokenBalanceKey(st.Key()):
			j, err := bs.handleTokenBalanceState(st)
			if err != nil {
				return err
			}
			TokenBalanceModels = append(TokenBalanceModels, j...)
		default:
			continue
		}
	}

	bs.tokenModels = TokenModels
	bs.tokenBalanceModels = TokenBalanceModels

	return nil
}

func (bs *BlockSession) handleTokenState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if tokenDoc, err := NewTokenDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(tokenDoc),
		}, nil
	}
}

func (bs *BlockSession) handleTokenBalanceState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if tokenBalanceDoc, err := NewTokenBalanceDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(tokenBalanceDoc),
		}, nil
	}
}
