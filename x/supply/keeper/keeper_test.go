package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/x/supply/types"
)

func TestValidatePermissions(t *testing.T) {
	app, _ := createTestApp(false)

	err := app.SupplyKeeper().ValidatePermissions(multiPermAcc)
	require.NoError(t, err)

	err = app.SupplyKeeper().ValidatePermissions(randomPermAcc)
	require.NoError(t, err)

	// unregistered permissions
	otherAcc := types.NewEmptyModuleAccount("other", "other")
	err = app.SupplyKeeper().ValidatePermissions(otherAcc)
	require.Error(t, err)
}
