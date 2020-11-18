package types

import (
	"fmt"

	ser "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	KuCodeSpace = "Kuchain"
)

func errorCode(root, sub uint32) uint32 {
	return root*10000 + sub
}

const (
	nameErrorCodeRoot = iota + 1
	kuMsgErrRoot
	authErrorCodeRoot
	txErrorCodeRoot
)

var (
	ErrNameParseTooLen = ser.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 1), "ErrNameParseTooLen")
	ErrNameNilString   = ser.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 2), "ErrNameNilString")
	ErrNameCharError   = ser.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 3), "ErrNameCharError")
	ErrNameStrNoValid  = ser.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 4), "ErrNameStrNoValid")
)

var (
	ErrKuMsgAuthCountTooLarge     = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 1), fmt.Sprintf("KuMsg Auth Count should less then %d", KuMsgMaxAuth))
	ErrKuMsgMissingFrom           = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 2), "KuMsg missing from accountID")
	ErrKuMsgMissingTo             = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 3), "KuMsg missing to accountID")
	ErrKuMsgMissingRouter         = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 4), "KuMsg missing router name")
	ErrKuMsgMissingType           = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 5), "KuMsg missing type name")
	ErrKuMsgMissingAuth           = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 6), "KuMsg missing auth for msg")
	ErrKuMsgDataTooLarge          = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 7), fmt.Sprintf("KuMsg msg data should <= %d", KuMsgMaxDataLen))
	ErrKuMsgDataUnmarshal         = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 8), "KuMsg msg data unmarshal error")
	ErrKuMsgDataSameAccount       = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 9), "KuMsg msg same account error")
	ErrKuMsgDataNotFindAccount    = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 10), "KuMsg msg can not find account error")
	ErrKuMsgAccountIDNil          = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 11), "KuMsg msg account id should not be nil")
	ErrKuMsgInconsistentAmount    = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 12), "KuMsg msg amount and data amount are inconsistent")
	ErrKuMSgNameEmpty             = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 13), "KuMsg msg name is empty")
	ErrKuMsgCoinsHasNegative      = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 14), "KuMsg coins has negative amount")
	ErrKuMsgFromNotEqual          = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 15), "KuMsg from not equal")
	ErrKuMsgToNotEqual            = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 16), "KuMsg to not equal")
	ErrKuMsgAmountNotEqual        = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 17), "KuMsg amount not equal")
	ErrKuMsgSpenderShouldNotEqual = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 18), "KuMsg apporve id and spender should not equal")
	ErrKuMsgTransferError         = ser.Register(KuCodeSpace, errorCode(kuMsgErrRoot, 19), "KuMsg transfer error")
)

var (
	ErrMissingAuth    = ser.Register(KuCodeSpace, errorCode(authErrorCodeRoot, 1), "Msg missing auth required")
	ErrTransfNoEnough = ser.Register(KuCodeSpace, errorCode(authErrorCodeRoot, 2), "Msg coin transf no enough")
	ErrTransfNotTo    = ser.Register(KuCodeSpace, errorCode(authErrorCodeRoot, 3), "Msg to account error")
)

var (
	ErrGasOverflow     = ser.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 1), "tx invalid gas supplied")
	ErrInsufficientFee = ser.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 2), "tx invalid fee amount provided")
	ErrNoSignatures    = ser.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 3), "tx no signers")
	ErrUnauthorized    = ser.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 4), "tx wrong number of signers")
	ErrTxDecode        = ser.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 5), "tx error decoding")
)
