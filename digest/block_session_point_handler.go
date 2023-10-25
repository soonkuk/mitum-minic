package digest

import (
	"github.com/ProtoconNet/mitum-point/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) preparePoint() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var PointModels []mongo.WriteModel
	var PointBalanceModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]

		switch {
		case state.IsStateDesignKey(st.Key()):
			j, err := bs.handlePointState(st)
			if err != nil {
				return err
			}
			PointModels = append(PointModels, j...)
		case state.IsStatePointBalanceKey(st.Key()):
			j, err := bs.handlePointBalanceState(st)
			if err != nil {
				return err
			}
			PointBalanceModels = append(PointBalanceModels, j...)
		default:
			continue
		}
	}

	bs.pointModels = PointModels
	bs.pointBalanceModels = PointBalanceModels

	return nil
}

func (bs *BlockSession) handlePointState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if pointDoc, err := NewPointDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(pointDoc),
		}, nil
	}
}

func (bs *BlockSession) handlePointBalanceState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if pointBalanceDoc, err := NewPointBalanceDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(pointBalanceDoc),
		}, nil
	}
}
