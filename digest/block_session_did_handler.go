package digest

import (
	"github.com/ProtoconNet/mitum-credential/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareDID() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var didModels []mongo.WriteModel
	var didCredentialModels []mongo.WriteModel
	var didHolderDIDModels []mongo.WriteModel
	var didTemplateModels []mongo.WriteModel

	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case state.IsStateDesignKey(st.Key()):
			j, err := bs.handleDIDIssuerState(st)
			if err != nil {
				return err
			}
			didModels = append(didModels, j...)
		case state.IsStateCredentialKey(st.Key()):
			j, err := bs.handleCredentialState(st)
			if err != nil {
				return err
			}
			didCredentialModels = append(didCredentialModels, j...)
		case state.IsStateHolderDIDKey(st.Key()):
			j, err := bs.handleHolderDIDState(st)
			if err != nil {
				return err
			}
			didHolderDIDModels = append(didHolderDIDModels, j...)
		case state.IsStateTemplateKey(st.Key()):
			j, err := bs.handleTemplateState(st)
			if err != nil {
				return err
			}
			didTemplateModels = append(didTemplateModels, j...)
		default:
			continue
		}
	}

	bs.didIssuerModels = didModels
	bs.didCredentialModels = didCredentialModels
	bs.didHolderDIDModels = didHolderDIDModels
	bs.didTemplateModels = didTemplateModels

	return nil
}

func (bs *BlockSession) handleDIDIssuerState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if issuerDoc, err := NewIssuerDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(issuerDoc),
		}, nil
	}
}

func (bs *BlockSession) handleCredentialState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if credentialDoc, err := NewCredentialDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(credentialDoc),
		}, nil
	}
}

func (bs *BlockSession) handleHolderDIDState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if holderDidDoc, err := NewHolderDIDDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(holderDidDoc),
		}, nil
	}
}

func (bs *BlockSession) handleTemplateState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if templateDoc, err := NewTemplateDoc(st, bs.st.DatabaseEncoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(templateDoc),
		}, nil
	}
}
