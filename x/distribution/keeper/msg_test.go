package keeper

import (
	"testing"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"

	"github.com/stretchr/testify/require"
)

// test ValidateBasic for MsgSetWithdrawAddress
func TestMsgSetWithdrawAddress(t *testing.T) {
	ctx, ak, _, _, _, _ := CreateTestInputDefault(t, false, 1000)

	emptyDel, _ := chainTypes.NewName("emptydel")
	emptyAcc := chainTypes.NewAccountIDFromName(emptyDel)

	tests := []struct {
		delegatorAddr types.AccountID
		withdrawAddr  types.AccountID
		expectPass    bool
	}{
		{Acc1, Acc2, true},
		{Acc3, Acc3, false},
		{emptyAcc, Acc4, false},
		{Acc5, emptyAcc, false},
		{emptyAcc, emptyAcc, false},
	}

	types.FindAcc = func(acc chainTypes.AccountID) bool {
		_, ok := acc.ToName()
		if ok {
			return ak.IsAccountExist(ctx, acc)
		}
		return false
	}

	for i, tc := range tests {
		Name, _ := tc.delegatorAddr.ToName()
		Auth, _ := ak.GetAuth(ctx, Name)

		msg := types.NewMsgSetWithdrawAccountID(Auth, tc.delegatorAddr, tc.withdrawAddr)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test index: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test index: %v", i)
		}
	}
}

// test ValidateBasic for MsgWithdrawDelegatorReward
func TestMsgWithdrawDelegatorReward(t *testing.T) {
	ctx, ak, _, _, _, _ := CreateTestInputDefault(t, false, 1000)

	emptyDel, _ := chainTypes.NewName("emptydel")
	emptyAcc := chainTypes.NewAccountIDFromName(emptyDel)

	tests := []struct {
		delegatorAddr AccountID
		validatorAddr AccountID
		expectPass    bool
	}{
		{Acc1, Acc2, true},
		{emptyAcc, Acc1, false},
		{Acc2, emptyAcc, false},
		{emptyAcc, emptyAcc, false},
	}
	types.FindAcc = func(acc chainTypes.AccountID) bool {
		_, ok := acc.ToName()
		if ok {
			return ak.IsAccountExist(ctx, acc)
		}
		return false
	}

	for i, tc := range tests {
		Name, _ := tc.delegatorAddr.ToName()
		Auth, _ := ak.GetAuth(ctx, Name)

		msg := types.NewMsgWithdrawDelegatorReward(Auth, tc.delegatorAddr, tc.validatorAddr)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test index: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test index: %v", i)
		}
	}
}

// test ValidateBasic for MsgWithdrawValidatorCommission
func TestMsgWithdrawValidatorCommission(t *testing.T) {
	ctx, ak, _, _, _, _ := CreateTestInputDefault(t, false, 1000)

	emptyDel, _ := chainTypes.NewName("emptydel")
	emptyAcc := chainTypes.NewAccountIDFromName(emptyDel)

	tests := []struct {
		validatorAddr AccountID
		expectPass    bool
	}{
		{Acc1, true},
		{emptyAcc, false},
	}

	types.FindAcc = func(acc chainTypes.AccountID) bool {
		_, ok := acc.ToName()
		if ok {
			return ak.IsAccountExist(ctx, acc)
		}
		return false
	}

	for i, tc := range tests {
		Name, _ := tc.validatorAddr.ToName()
		Auth, _ := ak.GetAuth(ctx, Name)

		msg := types.NewMsgWithdrawValidatorCommission(Auth, tc.validatorAddr)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test index: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test index: %v", i)
		}
	}
}
