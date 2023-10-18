package cmds

import (
	credentialcmds "github.com/ProtoconNet/mitum-credential/cmds"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	nftcmds "github.com/ProtoconNet/mitum-nft/v2/cmds"
	timestampcmds "github.com/ProtoconNet/mitum-timestamp/cmds"
	tokencmds "github.com/ProtoconNet/mitum-token/cmds"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	nftExtendedLen := currencyExtendedLen + len(nftcmds.AddedHinters)
	timestampExtendedLen := nftExtendedLen + len(timestampcmds.AddedHinters)
	credentialExtendedLen := timestampExtendedLen + len(credentialcmds.AddedHinters)
	allExtendedLen := credentialExtendedLen + len(tokencmds.AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:nftExtendedLen], nftcmds.AddedHinters)
	copy(Hinters[nftExtendedLen:timestampExtendedLen], timestampcmds.AddedHinters)
	copy(Hinters[timestampExtendedLen:credentialExtendedLen], credentialcmds.AddedHinters)
	copy(Hinters[credentialExtendedLen:], tokencmds.AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	nftSupportedExtendedLen := currencySupportedExtendedLen + len(nftcmds.AddedSupportedHinters)
	timestampSupportedExtendedLen := nftSupportedExtendedLen + len(timestampcmds.AddedSupportedHinters)
	credentialSupportedExtendedLen := timestampSupportedExtendedLen + len(credentialcmds.AddedSupportedHinters)
	allSupportedExtendedLen := credentialSupportedExtendedLen + len(tokencmds.AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:nftSupportedExtendedLen], nftcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[nftSupportedExtendedLen:timestampSupportedExtendedLen], timestampcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[timestampSupportedExtendedLen:credentialSupportedExtendedLen], credentialcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[credentialSupportedExtendedLen:], tokencmds.AddedSupportedHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
