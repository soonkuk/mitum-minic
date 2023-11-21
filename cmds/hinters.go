package cmds

import (
	credentialcmds "github.com/ProtoconNet/mitum-credential/cmds"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	daocmds "github.com/ProtoconNet/mitum-dao/cmds"
	nftcmds "github.com/ProtoconNet/mitum-nft/v2/cmds"
	pointcmds "github.com/ProtoconNet/mitum-point/cmds"
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
	tokenExtendedLen := credentialExtendedLen + len(tokencmds.AddedHinters)
	pointExtendedLen := tokenExtendedLen + len(pointcmds.AddedHinters)
	allExtendedLen := pointExtendedLen + len(daocmds.AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:nftExtendedLen], nftcmds.AddedHinters)
	copy(Hinters[nftExtendedLen:timestampExtendedLen], timestampcmds.AddedHinters)
	copy(Hinters[timestampExtendedLen:credentialExtendedLen], credentialcmds.AddedHinters)
	copy(Hinters[credentialExtendedLen:tokenExtendedLen], tokencmds.AddedHinters)
	copy(Hinters[tokenExtendedLen:pointExtendedLen], pointcmds.AddedHinters)
	copy(Hinters[pointExtendedLen:], daocmds.AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	nftSupportedExtendedLen := currencySupportedExtendedLen + len(nftcmds.AddedSupportedHinters)
	timestampSupportedExtendedLen := nftSupportedExtendedLen + len(timestampcmds.AddedSupportedHinters)
	credentialSupportedExtendedLen := timestampSupportedExtendedLen + len(credentialcmds.AddedSupportedHinters)
	tokenSupportedExtendedLen := credentialSupportedExtendedLen + len(tokencmds.AddedSupportedHinters)
	pointSupportedExtendedLen := tokenSupportedExtendedLen + len(pointcmds.AddedSupportedHinters)
	allSupportedExtendedLen := pointSupportedExtendedLen + len(daocmds.AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:nftSupportedExtendedLen], nftcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[nftSupportedExtendedLen:timestampSupportedExtendedLen], timestampcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[timestampSupportedExtendedLen:credentialSupportedExtendedLen], credentialcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[credentialSupportedExtendedLen:tokenSupportedExtendedLen], tokencmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[tokenSupportedExtendedLen:pointSupportedExtendedLen], pointcmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[pointSupportedExtendedLen:], daocmds.AddedSupportedHinters)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
