package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-minic-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var f currencycmds.ProposalOperationFactHintFunc = IsSupportedProposalOperationFactHintFunc
	pctx = context.WithValue(pctx, currencycmds.ProposalOperationFactHintContextKey, f)

	return pctx, nil
}
