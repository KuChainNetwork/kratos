package simapp

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
)

func ShouldErrIs(actual interface{}, expected ...interface{}) string {
	if len(expected) == 0 {
		return "`ShouldErrIs` must has a value to be check"
	}

	err1 := actual.(error)
	err2 := expected[0].(error)

	if errors.Is(err1, err2) {
		return ""
	}

	return fmt.Sprintf("Err %v is not be is %v!", actual, expected[0])
}

type bytesType interface {
	Bytes() []byte
}

func ShouldEq(actual interface{}, expected ...interface{}) string {
	if len(expected) == 0 {
		return "ShouldEq should between two values!"
	}

	if lName, lOk := actual.(types.Name); lOk {
		if rName, rOk := expected[0].(types.Name); rOk && lName.Eq(rName) {
			return ""
		}
	}

	if l, lOk := actual.(types.AccountID); lOk {
		if r, rOk := expected[0].(types.AccountID); rOk && l.Eq(r) {
			return ""
		}
	}

	if l, lOk := actual.(types.AccAddress); lOk {
		if r, rOk := expected[0].(types.AccAddress); rOk && l.Equals(r) {
			return ""
		}
	}

	if l, lOk := actual.(types.Coins); lOk {
		if r, rOk := expected[0].(types.Coins); rOk && l.IsEqual(r) {
			return ""
		}
	}

	if l, lOk := actual.(types.Coin); lOk {
		if r, rOk := expected[0].(types.Coin); rOk && l.IsEqual(r) {
			return ""
		}
	}

	// check is bytes equal
	l, lOk := actual.(bytesType)
	r, rOk := expected[0].(bytesType)
	if lOk && rOk && bytes.Equal(l.Bytes(), r.Bytes()) {
		return ""
	}

	return fmt.Sprintf("%v should equal %v", actual, expected[0])

}
