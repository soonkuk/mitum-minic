package processor

import (
	"fmt"
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-timestamp/operation/timestamp"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender                currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency              currencytypes.DuplicationType = "currency"
	DuplicationTypeContractCredential    currencytypes.DuplicationType = "contract-credential"
	DuplicationTypeContractTimeStamp     currencytypes.DuplicationType = "contract-timestamp"
	DuplicationTypeContractNFTCollection currencytypes.DuplicationType = "contract-collection"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op base.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var duplicationTypeSenderID string
	var duplicationTypeCurrencyID string
	var duplicationTypeContractCredentialID string
	var duplicationTypeContractTimestampID string
	var duplicationTypeContractNFTCollectionID string
	var newAddresses []base.Address

	switch t := op.(type) {
	case currency.CreateAccounts:
		fact, ok := t.Fact().(currency.CreateAccountsFact)
		if !ok {
			return errors.Errorf("expected CreateAccountsFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case currency.KeyUpdater:
		fact, ok := t.Fact().(currency.KeyUpdaterFact)
		if !ok {
			return errors.Errorf("expected KeyUpdaterFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Target().String()
	case currency.Transfers:
		fact, ok := t.Fact().(currency.TransfersFact)
		if !ok {
			return errors.Errorf("expected TransfersFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case currency.CurrencyRegister:
		fact, ok := t.Fact().(currency.CurrencyRegisterFact)
		if !ok {
			return errors.Errorf("expected CurrencyRegisterFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = fact.Currency().Currency().String()
	case currency.CurrencyPolicyUpdater:
		fact, ok := t.Fact().(currency.CurrencyPolicyUpdaterFact)
		if !ok {
			return errors.Errorf("expected CurrencyPolicyUpdaterFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Currency().String()
	case currency.SuffrageInflation:
	case extensioncurrency.CreateContractAccounts:
		fact, ok := t.Fact().(extensioncurrency.CreateContractAccountsFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountsFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case extensioncurrency.Withdraws:
		fact, ok := t.Fact().(extensioncurrency.WithdrawsFact)
		if !ok {
			return errors.Errorf("expected WithdrawsFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.CollectionRegister:
		fact, ok := t.Fact().(nft.CollectionRegisterFact)
		if !ok {
			return errors.Errorf("expected CollectionRegisterFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
		duplicationTypeContractNFTCollectionID = fact.Contract().String() + "-" + fact.Collection().String()
	case nft.CollectionPolicyUpdater:
		fact, ok := t.Fact().(nft.CollectionPolicyUpdaterFact)
		if !ok {
			return errors.Errorf("expected CollectionPolicyUpdaterFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.Mint:
		fact, ok := t.Fact().(nft.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.NFTTransfer:
		fact, ok := t.Fact().(nft.NFTTransferFact)
		if !ok {
			return errors.Errorf("expected NFTTransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.Delegate:
		fact, ok := t.Fact().(nft.DelegateFact)
		if !ok {
			return errors.Errorf("expected DelegateFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.Approve:
		fact, ok := t.Fact().(nft.ApproveFact)
		if !ok {
			return errors.Errorf("expected ApproveFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.NFTSign:
		fact, ok := t.Fact().(nft.NFTSignFact)
		if !ok {
			return errors.Errorf("expected NFTSignFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case timestamp.ServiceRegister:
		fact, ok := t.Fact().(timestamp.ServiceRegisterFact)
		if !ok {
			return errors.Errorf("expected ServiceRegisterFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
		duplicationTypeContractTimestampID = fact.Target().String() + "-" + fact.Service().String()
	case timestamp.Append:
		fact, ok := t.Fact().(timestamp.AppendFact)
		if !ok {
			return errors.Errorf("expected AppendFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case credential.CreateCredentialService:
		fact, ok := t.Fact().(credential.CreateCredentialServiceFact)
		if !ok {
			return errors.Errorf("expected CreateCredentialServiceFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
		duplicationTypeContractCredentialID = fact.Contract().String() + "-" + fact.CredentialServiceID().String()
	case credential.AddTemplate:
		fact, ok := t.Fact().(credential.AddTemplateFact)
		if !ok {
			return errors.Errorf("expected AddTemplateFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case credential.AssignCredentials:
		fact, ok := t.Fact().(credential.AssignCredentialsFact)
		if !ok {
			return errors.Errorf("expected AssignCredentialsFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case credential.RevokeCredentials:
		fact, ok := t.Fact().(credential.AssignCredentialsFact)
		if !ok {
			return errors.Errorf("expected RevokeCredentials, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	default:
		return nil
	}

	if len(duplicationTypeSenderID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeSenderID]; found {
			fmt.Println(">>>>>>>>>>>> duplication sender")
			return errors.Errorf("proposal cannot have duplicate sender, %v", duplicationTypeSenderID)
		}

		opr.Duplicated[duplicationTypeSenderID] = DuplicationTypeSender
	}
	if len(duplicationTypeCurrencyID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeSenderID]; found {
			return errors.Errorf(
				"cannot register duplicate currency id, %v within a proposal",
				duplicationTypeCurrencyID,
			)
		}

		opr.Duplicated[duplicationTypeSenderID] = DuplicationTypeCurrency
	}
	if len(duplicationTypeContractNFTCollectionID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContractNFTCollectionID]; found {
			return errors.Errorf(
				"cannot register a duplicate combination of contract-collection, %v within a proposal",
				duplicationTypeContractNFTCollectionID,
			)
		}

		opr.Duplicated[duplicationTypeContractNFTCollectionID] = DuplicationTypeContractNFTCollection
	}
	if len(duplicationTypeContractTimestampID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContractTimestampID]; found {
			return errors.Errorf(
				"cannot register a duplicate combination of contract-timestamp, %v within a proposal",
				duplicationTypeContractTimestampID,
			)
		}

		opr.Duplicated[duplicationTypeContractTimestampID] = DuplicationTypeContractTimeStamp
	}
	if len(duplicationTypeContractCredentialID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContractCredentialID]; found {
			return errors.Errorf(
				"cannot register a duplicate combination of contract-credential, %v within a proposal",
				duplicationTypeContractCredentialID,
			)
		}

		opr.Duplicated[duplicationTypeContractCredentialID] = DuplicationTypeContractCredential
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}
