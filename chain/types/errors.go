package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	KuCodeSpace = "Kuchain"
)

func errorCode(root, sub uint32) uint32 {
	return root*10000 + sub
}

const (
	nameErrorCodeRoot = iota + 1
	kuMsgErrorCodeRoot
	authErrorCodeRoot
	txErrorCodeRoot
)

var (
	ErrNameParseTooLen = sdkerrors.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 1), "ErrNameParseTooLen")
	ErrNameNilString   = sdkerrors.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 2), "ErrNameNilString")
	ErrNameCharError   = sdkerrors.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 3), "ErrNameCharError")
	ErrNameStrNoValid  = sdkerrors.Register(KuCodeSpace, errorCode(nameErrorCodeRoot, 4), "ErrNameStrNoValid")
)

var (
	ErrKuMsgAuthCountTooLarge     = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 1), fmt.Sprintf("KuMsg Auth Count should less then %d", KuMsgMaxAuth))
	ErrKuMsgMissingFrom           = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 2), "KuMsg missing from accountID")
	ErrKuMsgMissingTo             = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 3), "KuMsg missing to accountID")
	ErrKuMsgMissingRouter         = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 4), "KuMsg missing router name")
	ErrKuMsgMissingType           = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 5), "KuMsg missing type name")
	ErrKuMsgMissingAuth           = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 6), "KuMsg missing auth for msg")
	ErrKuMsgDataTooLarge          = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 7), fmt.Sprintf("KuMsg msg data should <= %d", KuMsgMaxDataLen))
	ErrKuMsgDataUnmarshal         = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 8), "KuMsg msg data unmarshal error")
	ErrKuMsgDataSameAccount       = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 9), "KuMsg msg same account error")
	ErrKuMsgDataNotFindAccount    = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 10), "KuMsg msg can not find account error")
	ErrKuMsgAccountIDNil          = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 11), "KuMsg msg account id should not be nil")
	ErrKuMsgInconsistentAmount    = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 12), "KuMsg msg amount and data amount are inconsistent")
	ErrKuMSgNameEmpty             = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 13), "KuMsg msg name is empty")
	ErrKuMsgCoinsHasNegative      = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 14), "KuMsg coins has negative amount")
	ErrKuMsgFromNotEqual          = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 15), "KuMsg from not equal")
	ErrKuMsgToNotEqual            = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 16), "KuMsg to not equal")
	ErrKuMsgAmountNotEqual        = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 17), "KuMsg amount not equal")
	ErrKuMsgSpenderShouldNotEqual = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 18), "KuMsg apporve id and spender should not equal")
	ErrKuMsgTransferError         = sdkerrors.Register(KuCodeSpace, errorCode(kuMsgErrorCodeRoot, 19), "KuMsg transfer error")
)

var (
	ErrMissingAuth    = sdkerrors.Register(KuCodeSpace, errorCode(authErrorCodeRoot, 1), "Msg missing auth required")
	ErrTransfNoEnough = sdkerrors.Register(KuCodeSpace, errorCode(authErrorCodeRoot, 2), "Msg coin transf no enough")
	ErrTransfNotTo    = sdkerrors.Register(KuCodeSpace, errorCode(authErrorCodeRoot, 3), "Msg to account error")
)

var (
	ErrGasOverflow     = sdkerrors.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 1), "tx invalid gas supplied")
	ErrInsufficientFee = sdkerrors.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 2), "tx invalid fee amount provided")
	ErrNoSignatures    = sdkerrors.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 3), "tx no signers")
	ErrUnauthorized    = sdkerrors.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 4), "tx wrong number of signers")
	ErrTxDecode        = sdkerrors.Register(KuCodeSpace, errorCode(txErrorCodeRoot, 5), "tx error decoding")
)
