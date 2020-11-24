package simapp

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/msg"
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/smartystreets/goconvey/convey"
)

type MsgInterface interface {
	types.Msg
	SetData(data []byte)
}

type noMsg struct {
	types.KuMsg
}

func (noMsg) ValidateBasic() error {
	return nil
}

type noExistMsgData struct {
	I int `json:"i" yaml:"i"`
}

func (noExistMsgData) Type() types.Name { return types.MustName("noexist") }

func (m noExistMsgData) Sender() types.AccountID {
	return types.MustAccountID("nonono")
}

func TestHandlerDataErr(t *testing.T, handler msg.Handler, msgs ...MsgInterface) {
	Convey("test handler data error", t, func() {
		Convey("test unknown req", func() {
			// check unknown req
			_, err := handler(types.Context{}, noMsg{})
			So(err, ShouldErrIs, sdkerrors.ErrUnknownRequest)
		})

		Convey("test unmarshal req", func() {
			for i, msg2Test := range msgs {
				// check unknown req
				dataByte, err := ModuleCdc.MarshalBinaryLengthPrefixed(noExistMsgData{I: i})
				if err != nil {
					panic(err)
				}

				msg2Test.SetData(dataByte)
				_, err = handler(types.Context{}, msg2Test)
				So(err, ShouldNotBeNil)
			}
		})
	})
}
