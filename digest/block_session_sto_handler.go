package digest

import (
	ststo "github.com/ProtoconNet/mitum-sto/state/sto"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareSTO() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var stoDesignModels []mongo.WriteModel
	var stoHolderPartitionsModels []mongo.WriteModel
	var stoHolderPartitionBalanceModels []mongo.WriteModel
	var stoHolderPartitionOperatorsModels []mongo.WriteModel
	var stoPartitionBalanceModels []mongo.WriteModel
	//var stoPartitionControllersModels []mongo.WriteModel
	var stoOperatorHoldersModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case ststo.IsStateDesignKey(st.Key()):
			j, err := bs.handleSTODesignState(st)
			if err != nil {
				return err
			}
			stoDesignModels = append(stoDesignModels, j...)
		case ststo.IsStateTokenHolderPartitionsKey(st.Key()):
			j, err := bs.handleSTOHolderPartitionsState(st)
			if err != nil {
				return err
			}
			stoHolderPartitionsModels = append(stoHolderPartitionsModels, j...)
		case ststo.IsStateTokenHolderPartitionBalanceKey(st.Key()):
			j, err := bs.handleSTOHolderPartitionBalanceState(st)
			if err != nil {
				return err
			}
			stoHolderPartitionBalanceModels = append(stoHolderPartitionBalanceModels, j...)
		case ststo.IsStateTokenHolderPartitionOperatorsKey(st.Key()):
			j, err := bs.handleSTOHolderPartitionOperatorsState(st)
			if err != nil {
				return err
			}
			stoHolderPartitionOperatorsModels = append(stoHolderPartitionOperatorsModels, j...)
		case ststo.IsStatePartitionBalanceKey(st.Key()):
			j, err := bs.handleSTOPartitionBalanceState(st)
			if err != nil {
				return err
			}
			stoPartitionBalanceModels = append(stoPartitionBalanceModels, j...)
		//case stostate.IsStatePartitionControllersKey(st.Key()):
		//	j, err := bs.handlePartitionControllersState(st)
		//	if err != nil {
		//		return err
		//	}
		//	stoPartitionControllersModels = append(stoPartitionControllersModels, j...)
		case ststo.IsStateOperatorTokenHoldersKey(st.Key()):
			j, err := bs.handleSTOperatorHoldersState(st)
			if err != nil {
				return err
			}
			stoOperatorHoldersModels = append(stoOperatorHoldersModels, j...)
		default:
			continue
		}
	}

	bs.stoDesignModels = stoDesignModels
	bs.stoHolderPartitionsModels = stoHolderPartitionsModels
	bs.stoHolderPartitionBalanceModels = stoHolderPartitionBalanceModels
	bs.stoHolderPartitionOperatorsModels = stoHolderPartitionOperatorsModels
	bs.stoPartitionBalanceModels = stoPartitionBalanceModels
	//bs.stoPartitionControllersModels = stoPartitionControllersModels
	bs.stoOperatorHoldersModels = stoOperatorHoldersModels

	return nil
}

func (bs *BlockSession) handleSTODesignState(st base.State) ([]mongo.WriteModel, error) {
	if doc, err := NewSTODesignDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOHolderPartitionsState(st base.State) ([]mongo.WriteModel, error) {
	if doc, err := NewSTOHolderPartitionsDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOHolderPartitionBalanceState(st base.State) ([]mongo.WriteModel, error) {
	if doc, err := NewSTOHolderPartitionBalanceDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOHolderPartitionOperatorsState(st base.State) ([]mongo.WriteModel, error) {
	if doc, err := NewSTOHolderPartitionOperatorsDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOPartitionBalanceState(st base.State) ([]mongo.WriteModel, error) {
	if doc, err := NewSTOPartitionBalanceDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}

func (bs *BlockSession) handleSTOperatorHoldersState(st base.State) ([]mongo.WriteModel, error) {
	if doc, err := NewSTOOperatorHoldersDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(doc),
		}, nil
	}
}
