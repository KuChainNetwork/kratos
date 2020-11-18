package types

import (
	"bytes"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/math"
)

const (
	accountIDTypeNil byte = iota
	accountIDTypeName
	accountIDTypeAccAddress
	accountIDTypeIdxLen
)

const (
	AccIDStoreKeyLen = sdk.AddrLen + 1
)

// AccountID will imp the cosmos Address interface
var _ sdk.Address = &AccountID{}

// AccountID a id for the entity which can own asset, now the accountID will be a AccAddress or a name for a account
type AccountID struct {
	Value []byte `json:"value,omitempty" yaml:"value"`
}

// AccountIDes defines a repeated set of AccountIDes.
type AccountIDes struct {
	Addresses []AccountID `json:"addresses"`
}

// NewAccountIDFromName create AccountID from Name
func NewAccountIDFromName(n Name) (res AccountID) {
	res.Value = make([]byte, AccIDStoreKeyLen)
	copy(res.Value, n.Value)
	return res
}

// NewAccountIDFromByte create AccountID from byte
func NewAccountIDFromByte(bytes []byte) AccountID {
	bytesNew := make([]byte, math.MaxInt(len(bytes), AccIDStoreKeyLen))
	copy(bytesNew, bytes)

	return AccountID{
		Value: bytesNew,
	}
}

// NewAccountIDFromAccAdd create AccountID from AccAddress
func NewAccountIDFromAccAdd(add sdk.AccAddress) AccountID {
	bs := make([]byte, 0, len(add)+1)
	bs = append(bs, accountIDTypeAccAddress)
	bs = append(bs, add...)
	return NewAccountIDFromByte(bs)
}

// NewAccountIDFromAdd create AccountID from AccAddress
func NewAccountIDFromAdd(add sdk.Address) AccountID {
	bs := make([]byte, 0, len(add.Bytes())+1)
	bs = append(bs, accountIDTypeAccAddress)
	bs = append(bs, add.Bytes()...)
	return NewAccountIDFromByte(bs)
}

// NewAccountIDFromValAdd create AccountID from AccAddress
func NewAccountIDFromValAdd(add sdk.ValAddress) AccountID {
	return NewAccountIDFromAccAdd(sdk.AccAddress(add))
}

// NewAccountIDFromConsAdd create AccountID from AccAddress
func NewAccountIDFromConsAdd(add sdk.ConsAddress) AccountID {
	return NewAccountIDFromAccAdd(sdk.AccAddress(add))
}

// EmptyAccountID return a empty accountID
func EmptyAccountID() AccountID {
	return NewAccountIDFromByte([]byte{accountIDTypeNil})
}

// NewAccountIDFromStr new accountID from string
func NewAccountIDFromStr(str string) (AccountID, error) {
	if str == "" {
		return EmptyAccountID(), nil
	}

	if len(str) <= NameStrLenMax {
		n, err := NewName(str)
		if err != nil {
			return AccountID{}, err
		}

		return NewAccountIDFromName(n), nil
	}

	accAdd, err := sdk.AccAddressFromBech32(str)
	if err != nil {
		return AccountID{}, err
	}

	return NewAccountIDFromAccAdd(accAdd), nil
}

// MustAccountID new accountID from string, if error then panic
func MustAccountID(str string) AccountID {
	res, err := NewAccountIDFromStr(str)
	if err != nil {
		panic(errors.Wrapf(err, "must accountID %s", str))
	}

	return res
}

// Equals if is same byte
func (a AccountID) Equals(o sdk.Address) bool {
	if a.Empty() && o.Empty() {
		return true
	}

	if acc, ok := o.(sdk.AccAddress); ok {
		acc2, ok := a.ToAccAddress()
		return ok && acc2.Equals(acc)
	}

	return bytes.Equal(a.Bytes(), o.Bytes())
}

// Equals if is same byte
func (a AccountID) Eq(o AccountID) bool {
	if a.Empty() && o.Empty() {
		return true
	}

	return bytes.Equal(a.Bytes(), o.Bytes())
}

// Equals if is same byte
func (a AccountID) Equal(o *AccountID) bool {
	if a.Empty() && o.Empty() {
		return true
	}

	return bytes.Equal(a.Bytes(), o.Bytes())
}

// Empty return is AccountID is empty, if AccountID is a name, return is a empty name
func (a AccountID) Empty() bool {
	if len(a.Value) == 0 {
		return true
	}

	if a.Value[0] == accountIDTypeNil {
		return true
	}

	if addr, ok := a.ToAddress(); ok {
		return addr.Empty()
	}

	if name, ok := a.ToName(); ok {
		return name.Empty()
	}

	return true
}

// MarshalJSON imp Address interface
func (a AccountID) MarshalJSON() ([]byte, error) {
	if a.Empty() {
		return json.Marshal("")
	}

	if addr, ok := a.ToAddress(); ok {
		return addr.MarshalJSON()
	}

	if name, ok := a.ToName(); ok {
		return name.MarshalJSON()
	}

	return []byte{}, nil
}

// UnmarshalJSON unmarshal from JSON assuming Bech32 encoding.
func (a *AccountID) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	if len(s) == 0 {
		*a = EmptyAccountID()
		return nil
	}

	if len(s) <= (NameStrLenMax + 2) {
		// must a name
		name, err := NewName(s)
		if err != nil {
			return err
		}
		*a = NewAccountIDFromName(name)
	} else {
		add := AccAddress{}

		if err := add.UnmarshalJSON(data); err != nil {
			return err
		}

		*a = NewAccountIDFromAccAdd(add)
	}

	return nil
}

// MarshalYAML imp to yaml
func (a AccountID) MarshalYAML() (interface{}, error) {
	if a.Empty() {
		return "", nil
	}

	if addr, ok := a.ToAddress(); ok {
		return addr.String(), nil
	}

	if name, ok := a.ToName(); ok {
		return name.MarshalYAML()
	}

	return []byte{}, nil
}

// Bytes imp Address interface
func (a AccountID) Bytes() []byte {
	if a.Empty() {
		return []byte{}
	}

	return a.StoreKey()
}

// String imp Address interface
func (a AccountID) String() string {
	if a.Empty() {
		return ""
	}

	if addr, ok := a.ToAddress(); ok {
		return addr.String()
	}

	if name, ok := a.ToName(); ok {
		return name.String()
	}

	return ""
}

func (a AccountID) Marshal() ([]byte, error) {
	return a.Bytes(), nil
}

// Format imp Address interface
func (a AccountID) Format(s fmt.State, verb rune) {
	if a.Empty() {
		s.Write([]byte{})
	}

	if addr, ok := a.ToAddress(); ok {
		addr.Format(s, verb)
	}

	if name, ok := a.ToName(); ok {
		name.Format(s, verb)
	}
}

// ToAddress if a is a address return address( a AccAddress ) and true
func (a AccountID) ToAddress() (sdk.Address, bool) {
	if len(a.Value) == 0 {
		return nil, false
	}

	return a.ToAccAddress()
}

// ToAccAddress if a is a account address return AccAddress and true
func (a AccountID) ToAccAddress() (sdk.AccAddress, bool) {
	if a.Value[0] == accountIDTypeAccAddress {
		cb := make([]byte, len(a.Value)-1)
		copy(cb, a.Value[1:])
		return cb, true
	}

	return nil, false
}

// ToName if `a` is a name return name and true
func (a AccountID) ToName() (Name, bool) {
	if len(a.Value) == 0 {
		return Name{}, false
	}

	if a.Value[0] == accountIDTypeName {
		return NewNameFromBytes(a.Value), true
	}

	return Name{}, false
}

// MustAddress if a is a address return address( a AccAddress ) and true
func (a AccountID) MustAddress() sdk.Address {
	res, ok := a.ToAddress()
	if !ok {
		panic(fmt.Errorf("accountID no Address %s", a))
	}
	return res
}

// MustAccAddress if a is a account address return AccAddress and true
func (a AccountID) MustAccAddress() sdk.AccAddress {
	res, ok := a.ToAccAddress()
	if !ok {
		panic(fmt.Errorf("accountID no AccAddress %s", a))
	}
	return res
}

// MustName if `a` is a name return name and true
func (a AccountID) MustName() Name {
	res, ok := a.ToName()
	if !ok {
		panic(fmt.Errorf("accountID no Name %s", a))
	}
	return res
}

// NewAccountIDFromStoreKey creates a new accountID from store key
func NewAccountIDFromStoreKey(val []byte) AccountID {
	v := val[1:]
	if len(v) != AccIDStoreKeyLen {
		panic("unexpected AccountID store key length")
	}

	if v[0] >= accountIDTypeIdxLen {
		panic("unexpected AccountID prefix store key length")
	}

	return AccountID{
		Value: v,
	}
}

// StoreKey get bytes for store key
func (a AccountID) StoreKey() []byte {
	if a.Empty() {
		return []byte{}
	}

	bytes := make([]byte, math.MaxInt(len(a.Value), AccIDStoreKeyLen))
	copy(bytes, a.Value)
	return bytes
}

// FixedKey get bytes for store key fixed
func (a AccountID) FixedKey() (res [AccIDStoreKeyLen]byte) {
	copy(res[:], a.Value)
	return
}
