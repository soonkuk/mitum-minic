package digest

import (
	timestampservice "github.com/ProtoconNet/mitum-timestamp/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareTimeStamps() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var timestampModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case timestampservice.IsStateServiceDesignKey(st.Key()):
			j, err := bs.handleTimeStampServiceDesignState(st)
			if err != nil {
				return err
			}
			timestampModels = append(timestampModels, j...)
		case timestampservice.IsStateTimeStampItemKey(st.Key()):
			j, err := bs.handleTimeStampItemState(st)
			if err != nil {
				return err
			}
			timestampModels = append(timestampModels, j...)
		default:
			continue
		}
	}

	bs.timestampModels = timestampModels

	return nil
}

func (bs *BlockSession) handleTimeStampServiceDesignState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if serviceDesignDoc, err := NewTimeStampServiceDesignDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(serviceDesignDoc),
		}, nil
	}
}

func (bs *BlockSession) handleTimeStampItemState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if TimeStampItemDoc, err := NewTimeStampItemDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(TimeStampItemDoc),
		}, nil
	}
}
