package keeper

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (a AssetKeeper) PayFee(ctx sdk.Context, payer types.AccountID, fee types.Coins) error {
	ctx.Logger().Debug("pay fee", "payer", payer, "fee", fee)
	if err := a.CoinsToPower(ctx, payer, constants.GetFeeCollector(), fee); err != nil {
		return sdkerrors.Wrap(err, "pay fee")
	}

	return nil
}
