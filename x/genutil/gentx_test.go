package genutil_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/genutil"
	"github.com/KuChainNetwork/kuchain/x/staking"
	stakingtypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	crypKeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/cli"
)

var smbh = staking.AppModuleBasic{}
var cmd = &cobra.Command{}
var wallet = simapp.NewWallet()

func TestSetGenTxsInAppGenesisState(t *testing.T) {
	Convey("TestSetGenTxsInAppGenesisState", t, func() {
		home, cleanup := simapp.NewTestCaseDir(t)
		defer cleanup()
		viper.Set(cli.HomeFlag, home)
		viper.Set(flags.FlagKeyringBackend, "test")

		cfg, err := tcmd.ParseConfig()
		cdc := makeCodec()
		appGenesisState := make(map[string]json.RawMessage, 0)
		inBuf := bufio.NewReader(os.Stdin)
		auth := wallet.NewAccAddress()
		smbh.PrepareFlagsForTxCreateValidator(cfg, "moniker", "test_chain", wallet.PrivKey(auth).PubKey())
		txBldr := txutil.NewTxBuilderFromCLI(inBuf).WithTxEncoder(txutil.GetTxEncoder(cdc)).WithPayer("validator")
		cliCtx := txutil.NewKuCLICtxByBuf(cdc, inBuf)
		valAccountID := types.MustAccountID("validator")

		txBldr, msg, err := smbh.BuildCreateValidatorMsg(cliCtx, txBldr, valAccountID, auth)
		So(err, ShouldBeNil)

		txBldr, msgdelegator, err := smbh.BuildDelegateMsg(cliCtx, txBldr, auth, valAccountID, valAccountID)
		So(err, ShouldBeNil)

		stdSignMsg, err := txBldr.BuildSignMsg([]sdk.Msg{msg, msgdelegator})
		So(err, ShouldBeNil)

		stdTx := types.NewStdTx(stdSignMsg.Msg, stdSignMsg.Fee, nil, stdSignMsg.Memo)

		appGenesisState, err = genutil.SetGenTxsInAppGenesisState(cdc, appGenesisState, []txutil.StdTx{stdTx})
		So(err, ShouldBeNil)
	})
}

func setKey(name string, inBuf *bufio.Reader) types.AccAddress {
	viper.Set(flags.FlagKeyringBackend, "test")
	viper.Set(cli.OutputFlag, "json")

	addKeyCmd := keys.AddKeyCommand()
	err := addKeyCmd.RunE(cmd, []string{name})
	So(err, ShouldBeNil)

	keys.ShowKeysCmd()

	kb, err := crypKeys.NewKeyring(sdk.KeyringServiceName(),
		viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), cmd.InOrStdin())
	So(err, ShouldBeNil)

	key, err := kb.Get(name)
	So(err, ShouldBeNil)

	auth := key.GetAddress()
	So(err, ShouldBeNil)

	return auth
}

func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	stakingtypes.RegisterCodec(cdc)
	return cdc
}
