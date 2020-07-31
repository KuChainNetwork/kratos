package cli

import (
	"bytes"
	"github.com/KuChainNetwork/kuchain/chain/client/flags"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/genutil"
	stakingtypes "github.com/KuChainNetwork/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	abciServer "github.com/tendermint/tendermint/abci/server"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"io"
	"os"
	"testing"
	"time"
)

var testMbm = module.NewBasicManager(genutil.AppModuleBasic{})

func TestInitCmd(t *testing.T) {
	Convey("TestInitCmd", t, func() {
		defer setupClientHome(t)()
		home, cleanup := simapp.NewTestCaseDir(t)
		defer cleanup()
		viper.Set(cli.HomeFlag, home)

		logger := log.NewNopLogger()
		cfg, err := tcmd.ParseConfig()
		So(err, ShouldBeNil)

		ctx := server.NewContext(cfg, logger)
		cdc := makeCodec()
		cmd := InitCmd(ctx, cdc, testMbm, home)

		err = cmd.RunE(nil, []string{"kuchain-test"})
		So(err, ShouldBeNil)
	})
}

func setupClientHome(t *testing.T) func() {
	clientDir, cleanup := simapp.NewTestCaseDir(t)
	viper.Set(flagClientHome, clientDir)
	return cleanup
}

func TestEmptyState(t *testing.T) {
	Convey("TestEmptyState", t, func() {
		defer setupClientHome(t)()
		home, cleanup := simapp.NewTestCaseDir(t)
		defer cleanup()
		viper.Set(cli.HomeFlag, home)

		logger := log.NewNopLogger()
		cfg, err := tcmd.ParseConfig()
		So(err, ShouldBeNil)

		ctx := server.NewContext(cfg, logger)
		cdc := makeCodec()

		cmd := InitCmd(ctx, cdc, testMbm, home)
		err = cmd.RunE(nil, []string{"kuchain-test"})
		So(err, ShouldBeNil)

		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		cmd = server.ExportCmd(ctx, cdc, nil)

		err = cmd.RunE(nil, nil)
		So(err, ShouldBeNil)

		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		w.Close()
		os.Stdout = old
		out := <-outC

		So(out, ShouldContainSubstring, "genesis_time")
		So(out, ShouldContainSubstring, "chain_id")
		So(out, ShouldContainSubstring, "consensus_params")
		So(out, ShouldContainSubstring, "app_hash")
		So(out, ShouldContainSubstring, "app_state")
	})
}

func TestStartStandAlone(t *testing.T) {
	Convey("TestStartStandAlone", t, func() {
		home, cleanup := simapp.NewTestCaseDir(t)
		defer cleanup()
		viper.Set(cli.HomeFlag, home)
		defer setupClientHome(t)()

		logger := log.NewNopLogger()
		cfg, err := tcmd.ParseConfig()
		So(err, ShouldBeNil)
		ctx := server.NewContext(cfg, logger)
		cdc := makeCodec()
		initCmd := InitCmd(ctx, cdc, testMbm, home)
		err = initCmd.RunE(nil, []string{"kuchain-test"})
		So(err, ShouldBeNil)

		app, err := mock.NewApp(home, logger)
		So(err, ShouldBeNil)
		svrAddr, _, err := server.FreeTCPAddr()
		So(err, ShouldBeNil)
		svr, err := abciServer.NewServer(svrAddr, "socket", app)
		So(err, ShouldBeNil)
		svr.SetLogger(logger.With("module", "abci-server"))
		svr.Start()

		timer := time.NewTimer(time.Duration(2) * time.Second)
		select {
		case <-timer.C:
			svr.Stop()
		}
	})
}

func TestInitNodeValidatorFiles(t *testing.T) {
	Convey("TestInitNodeValidatorFiles", t, func() {
		home, cleanup := simapp.NewTestCaseDir(t)
		defer cleanup()
		viper.Set(cli.HomeFlag, home)
		viper.Set(flags.FlagName, "moniker")
		cfg, err := tcmd.ParseConfig()
		So(err, ShouldBeNil)
		nodeID, valPubKey, err := genutil.InitializeNodeValidatorFiles(cfg)
		So(err, ShouldBeNil)
		So(nodeID, ShouldNotEqual, "")
		So(len(valPubKey.Bytes()), ShouldNotEqual, 0)
	})
}

// custom tx codec
func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	stakingtypes.RegisterCodec(cdc)
	return cdc
}
