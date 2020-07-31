package types

import (
	"encoding/json"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/constants/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	yaml "gopkg.in/yaml.v2"
)

var (
	_ sdk.Tx = (*StdTx)(nil)

	maxGasWanted = uint64((1 << 63) - 1)
)

// StdTx is a standard way to wrap a Msg with Fee and Signatures.
// NOTE: the first signature is the fee payer (Signatures must not be nil).
type StdTx struct {
	Msgs       []sdk.Msg      `json:"msg" yaml:"msg"`
	Fee        StdFee         `json:"fee" yaml:"fee"`
	Signatures []StdSignature `json:"signatures" yaml:"signatures"`
	Memo       string         `json:"memo" yaml:"memo"`
}

func NewStdTx(msgs []sdk.Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

// GetMsgs returns the all the transaction's messages.
func (tx StdTx) GetMsgs() []sdk.Msg { return tx.Msgs }

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() error {
	stdSigs := tx.GetSignatures()

	if tx.Fee.Gas > maxGasWanted {
		return errors.Wrapf(ErrGasOverflow, "%d > %d", tx.Fee.Gas, maxGasWanted)
	}

	if tx.Fee.Amount.IsAnyNegative() {
		return errors.Wrapf(ErrInsufficientFee, "%s amount provided", tx.Fee.Amount)
	}

	if len(stdSigs) == 0 {
		return ErrNoSignatures
	}

	if len(stdSigs) != len(tx.GetSigners()) {
		return ErrUnauthorized
	}

	if len(tx.Memo) > keys.DefaultMaxMemoCharacters {
		return sdkerrors.ErrMemoTooLarge
	}

	return nil
}

// CountSubKeys counts the total number of keys for a multi-sig public key.
func CountSubKeys(pub crypto.PubKey) int {
	v, ok := pub.(multisig.PubKeyMultisigThreshold)
	if !ok {
		return 1
	}

	numKeys := 0
	for _, subkey := range v.PubKeys {
		numKeys += CountSubKeys(subkey)
	}

	return numKeys
}

// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigners method for each Msg
// in the order they appear in tx.GetMsgs().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigners() []sdk.AccAddress {
	seen := map[string]bool{}
	var signers []sdk.AccAddress
	for _, msg := range tx.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}
	return signers
}

// GetMemo returns the memo
func (tx StdTx) GetMemo() string { return tx.Memo }

// GetSignatures returns the signature of signers who signed the Msg.
// GetSignatures returns the signature of signers who signed the Msg.
// CONTRACT: Length returned is same as length of
// pubkeys returned from MsgKeySigners, and the order
// matches.
// CONTRACT: If the signature is missing (ie the Msg is
// invalid), then the corresponding signature is
// .Empty().
func (tx StdTx) GetSignatures() []StdSignature { return tx.Signatures }

func (tx StdTx) PrettifyJSON(cdc *codec.Codec) ([]byte, error) {
	alias := struct {
		Msgs       []json.RawMessage `json:"msg" yaml:"msg"`
		Fee        StdFee            `json:"fee" yaml:"fee"`
		Signatures []StdSignature    `json:"signatures" yaml:"signatures"`
		Memo       string            `json:"memo" yaml:"memo"`
	}{
		Fee:        tx.Fee,
		Signatures: tx.Signatures,
		Memo:       tx.Memo,
		Msgs:       make([]json.RawMessage, 0, len(tx.Msgs)),
	}

	for _, msg := range tx.Msgs {
		if msg, ok := msg.(Prettifier); ok {
			raw, err := msg.PrettifyJSON(cdc)
			if err != nil {
				return []byte{}, sdkerrors.Wrapf(err, "prettify json to msg")
			}

			alias.Msgs = append(alias.Msgs, json.RawMessage(raw))
		}
	}

	return cdc.MarshalJSON(alias)
}

// GetGas returns the Gas in StdFee
func (tx StdTx) GetGas() uint64 { return tx.Fee.Gas }

// GetFee returns the FeeAmount in StdFee
func (tx StdTx) GetFee() Coins { return tx.Fee.Amount }

func (tx StdTx) FeePayer() AccountID {
	if !tx.Fee.Payer.Empty() {
		return tx.Fee.Payer
	}

	auths := tx.GetSigners()
	if auths != nil {
		return NewAccountIDFromAccAdd(auths[0])
	}

	return NewAccountIDFromAccAdd(sdk.AccAddress{})
}

//__________________________________________________________

// StdFee includes the amount of coins paid in fees and the maximum
// gas to be used by the transaction. The ratio yields an effective "gasprice",
// which must be above some miminum to be accepted into the mempool.
type StdFee struct {
	Amount Coins     `json:"amount" yaml:"amount"`
	Gas    uint64    `json:"gas" yaml:"gas"`
	Payer  AccountID `json:"payer" yaml:"payer"`
}

// NewStdFee returns a new instance of StdFee
func NewStdFee(gas uint64, payer AccountID, amount Coins) StdFee {
	return StdFee{
		Amount: amount,
		Gas:    gas,
		Payer:  payer,
	}
}

// Bytes for signing later
func (fee StdFee) Bytes() []byte {
	// normalize. XXX
	// this is a sign of something ugly
	// (in the lcd_test, client side its null,
	// server side its [])
	if len(fee.Amount) == 0 {
		fee.Amount = NewCoins()
	}
	bz, err := ModuleCdc.MarshalJSON(fee) // TODO
	if err != nil {
		panic(err)
	}
	return bz
}

// GasPrices returns the gas prices for a StdFee.
//
// NOTE: The gas prices returned are not the true gas prices that were
// originally part of the submitted transaction because the fee is computed
// as fee = ceil(gasWanted * gasPrices).
func (fee StdFee) GasPrices() DecCoins {
	return NewDecCoinsFromCoins(fee.Amount...).QuoDec(NewDec(int64(fee.Gas)))
}

//__________________________________________________________

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
	AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	ChainID       string            `json:"chain_id" yaml:"chain_id"`
	Fee           json.RawMessage   `json:"fee" yaml:"fee"`
	Memo          string            `json:"memo" yaml:"memo"`
	Msg           []json.RawMessage `json:"msg" yaml:"msg"`
	Sequence      uint64            `json:"sequence" yaml:"sequence"`
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, accnum uint64, sequence uint64, fee StdFee, msgs []sdk.Msg, memo string) []byte {
	var msgsBytes []json.RawMessage
	for _, msg := range msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}
	bz, err := ModuleCdc.MarshalJSON(StdSignDoc{
		AccountNumber: accnum,
		ChainID:       chainID,
		Fee:           json.RawMessage(fee.Bytes()),
		Memo:          memo,
		Msg:           msgsBytes,
		Sequence:      sequence,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// StdSignature represents a sig
type StdSignature struct {
	crypto.PubKey `json:"pub_key" yaml:"pub_key"` // optional
	Signature     []byte                          `json:"signature" yaml:"signature"`
}

// DefaultTxDecoder logic for standard transaction decoding
func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var tx = StdTx{}

		if len(txBytes) == 0 {
			return nil, errors.Wrap(ErrTxDecode, "txBytes are empty")
		}

		// StdTx.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding transaction")
		}

		return tx, nil
	}
}

// DefaultTxEncoder logic for standard transaction encoding
func DefaultTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.MarshalBinaryLengthPrefixed(tx)
	}
}

// MarshalYAML returns the YAML representation of the signature.
func (ss StdSignature) MarshalYAML() (interface{}, error) {
	var (
		bz     []byte
		pubkey string
		err    error
	)

	if ss.PubKey != nil {
		pubkey, err = sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, ss.PubKey)
		if err != nil {
			return nil, err
		}
	}

	bz, err = yaml.Marshal(struct {
		PubKey    string
		Signature string
	}{
		PubKey:    pubkey,
		Signature: fmt.Sprintf("%s", ss.Signature),
	})
	if err != nil {
		return nil, err
	}

	return string(bz), err
}
