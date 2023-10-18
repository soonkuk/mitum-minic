package processor

import (
	"github.com/ProtoconNet/mitum-credential/operation/credential"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-timestamp/operation/timestamp"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender   currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency currencytypes.DuplicationType = "currency"
	DuplicationTypeContract currencytypes.DuplicationType = "contract"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op base.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	var duplicationTypeSenderID string
	var duplicationTypeCurrencyID string
	var duplicationTypeContract string
	var newAddresses []base.Address

	switch t := op.(type) {
	case currency.CreateAccount:
		fact, ok := t.Fact().(currency.CreateAccountFact)
		if !ok {
			return errors.Errorf("expected CreateAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case currency.UpdateKey:
		fact, ok := t.Fact().(currency.UpdateKeyFact)
		if !ok {
			return errors.Errorf("expected UpdateKeyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Target().String()
	case currency.Transfer:
		fact, ok := t.Fact().(currency.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case currency.RegisterCurrency:
		fact, ok := t.Fact().(currency.RegisterCurrencyFact)
		if !ok {
			return errors.Errorf("expected RegisterCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeCurrencyID = fact.Currency().Currency().String()
	case currency.UpdateCurrency:
		fact, ok := t.Fact().(currency.UpdateCurrencyFact)
		if !ok {
			return errors.Errorf("expected UpdateCurrencyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Currency().String()
	case currency.Mint:
	case extensioncurrency.CreateContractAccount:
		fact, ok := t.Fact().(extensioncurrency.CreateContractAccountFact)
		if !ok {
			return errors.Errorf("expected CreateContractAccountFact, not %T", t.Fact())
		}
		as, err := fact.Targets()
		if err != nil {
			return errors.Errorf("failed to get Addresses")
		}
		newAddresses = as
		duplicationTypeSenderID = fact.Sender().String()
	case extensioncurrency.Withdraw:
		fact, ok := t.Fact().(extensioncurrency.WithdrawFact)
		if !ok {
			return errors.Errorf("expected WithdrawFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.CreateCollection:
		fact, ok := t.Fact().(nft.CreateCollectionFact)
		if !ok {
			return errors.Errorf("expected CreateCollectionFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
		duplicationTypeContract = fact.Contract().String()
	case nft.UpdateCollectionPolicy:
		fact, ok := t.Fact().(nft.UpdateCollectionPolicyFact)
		if !ok {
			return errors.Errorf("expected UpdateCollectionPolicyFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.Mint:
		fact, ok := t.Fact().(nft.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case nft.Transfer:
		fact, ok := t.Fact().(nft.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
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
	case nft.Sign:
		fact, ok := t.Fact().(nft.SignFact)
		if !ok {
			return errors.Errorf("expected SignFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case timestamp.CreateService:
		fact, ok := t.Fact().(timestamp.CreateServiceFact)
		if !ok {
			return errors.Errorf("expected CreateServiceFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
		duplicationTypeContract = fact.Target().String()
	case timestamp.Append:
		fact, ok := t.Fact().(timestamp.AppendFact)
		if !ok {
			return errors.Errorf("expected AppendFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case credential.CreateService:
		fact, ok := t.Fact().(credential.CreateServiceFact)
		if !ok {
			return errors.Errorf("expected CreateServiceFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
		duplicationTypeContract = fact.Contract().String()
	case credential.AddTemplate:
		fact, ok := t.Fact().(credential.AddTemplateFact)
		if !ok {
			return errors.Errorf("expected AddTemplateFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case credential.Assign:
		fact, ok := t.Fact().(credential.AssignFact)
		if !ok {
			return errors.Errorf("expected AssignFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case credential.Revoke:
		fact, ok := t.Fact().(credential.RevokeFact)
		if !ok {
			return errors.Errorf("expected RevokeFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case token.Mint:
		fact, ok := t.Fact().(token.MintFact)
		if !ok {
			return errors.Errorf("expected MintFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case token.RegisterToken:
		fact, ok := t.Fact().(token.RegisterTokenFact)
		if !ok {
			return errors.Errorf("expected RegisterTokenFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
		duplicationTypeContract = fact.Contract().String()
	case token.Burn:
		fact, ok := t.Fact().(token.BurnFact)
		if !ok {
			return errors.Errorf("expected BurnFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case token.Approve:
		fact, ok := t.Fact().(token.ApproveFact)
		if !ok {
			return errors.Errorf("expected ApproveFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case token.Transfer:
		fact, ok := t.Fact().(token.TransferFact)
		if !ok {
			return errors.Errorf("expected TransferFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	case token.TransferFrom:
		fact, ok := t.Fact().(token.TransferFromFact)
		if !ok {
			return errors.Errorf("expected TransferFromFact, not %T", t.Fact())
		}
		duplicationTypeSenderID = fact.Sender().String()
	default:
		return nil
	}

	if len(duplicationTypeSenderID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeSenderID]; found {
			return errors.Errorf("proposal cannot have duplicate sender, %v", duplicationTypeSenderID)
		}

		opr.Duplicated[duplicationTypeSenderID] = DuplicationTypeSender
	}
	if len(duplicationTypeCurrencyID) > 0 {
		if _, found := opr.Duplicated[duplicationTypeCurrencyID]; found {
			return errors.Errorf(
				"cannot register duplicate currency id, %v within a proposal",
				duplicationTypeCurrencyID,
			)
		}

		opr.Duplicated[duplicationTypeCurrencyID] = DuplicationTypeCurrency
	}
	if len(duplicationTypeContract) > 0 {
		if _, found := opr.Duplicated[duplicationTypeContract]; found {
			return errors.Errorf(
				"cannot use a duplicated contract for registering in contract model , %v within a proposal",
				duplicationTypeSenderID,
			)
		}

		opr.Duplicated[duplicationTypeContract] = DuplicationTypeContract
	}

	if len(newAddresses) > 0 {
		if err := opr.CheckNewAddressDuplication(newAddresses); err != nil {
			return err
		}
	}

	return nil
}
