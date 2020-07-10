package params

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/params/external"
	"github.com/KuChainNetwork/kuchain/x/params/keeper"
	"github.com/KuChainNetwork/kuchain/x/params/types/proposal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParamChangeProposalHandler creates a new governance Handler for a ParamChangeProposal
func NewParamChangeProposalHandler(k keeper.Keeper) external.GovHandler {
	return func(ctx sdk.Context, content external.GovContent) error {
		switch c := content.(type) {
		case proposal.ParameterChangeProposal:
			return handleParameterChangeProposal(ctx, k, c)
		case *proposal.ParameterChangeProposal:
			return handleParameterChangeProposal(ctx, k, *c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized param proposal content type: %T", c)
		}
	}
}

func handleParameterChangeProposal(ctx sdk.Context, k keeper.Keeper, p proposal.ParameterChangeProposal) error {
	for _, c := range p.Changes {
		ss, ok := k.GetSubspace(c.Subspace)
		if !ok {
			return sdkerrors.Wrap(proposal.ErrUnknownSubspace, c.Subspace)
		}

		k.Logger(ctx).Info(
			fmt.Sprintf("attempt to set new parameter value; key: %s, value: %s", c.Key, c.Value),
		)

		if err := ss.Update(ctx, []byte(c.Key), []byte(c.Value)); err != nil {
			return sdkerrors.Wrapf(proposal.ErrSettingParameter, "key: %s, value: %s, err: %s", c.Key, c.Value, err.Error())
		}
	}

	return nil
}
