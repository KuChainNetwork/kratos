package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bytes"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// Implements Delegation interface
var _ exported.DelegationI = Delegation{}

// DVPair is struct that just has a delegator-validator pair with no other data.
// It is intended to be used as a marshalable pointer. For example, a DVPair can
// be used to construct the key to getting an UnbondingDelegation from state.
type DVPair struct {
	DelegatorAccount AccountID `json:"delegator_account" yaml:"delegator_account"`
	ValidatorAccount AccountID `json:"validator_account" yaml:"validator_account"`
}

// DVPairs defines an array of DVPair objects.
type DVPairs struct {
	Pairs []DVPair `json:"pairs" yaml:"pairs"`
}

// String implements the Stringer interface for a DVPair object.
func (dv DVPair) String() string {
	out, _ := yaml.Marshal(dv)
	return string(out)
}

// DVVTriplet is struct that just has a delegator-validator-validator triplet
// with no other data. It is intended to be used as a marshalable pointer. For
// example, a DVVTriplet can be used to construct the key to getting a
// Redelegation from state.
type DVVTriplet struct {
	DelegatorAccount    AccountID `json:"delegator_account" yaml:"delegator_account"`
	ValidatorSrcAccount AccountID `json:"validator_src_account" yaml:"validator_src_account"`
	ValidatorDstAccount AccountID `json:"validator_dst_account" yaml:"validator_dst_account"`
}

// DVVTriplets defines an array of DVVTriplet objects.
type DVVTriplets struct {
	Triplets []DVVTriplet `json:"triplets" yaml:"triplets"`
}

// String implements the Stringer interface for a DVVTriplet object.
func (dvv DVVTriplet) String() string {
	out, _ := yaml.Marshal(dvv)
	return string(out)
}

// Delegation represents the bond with tokens held by an account. It is
// owned by one delegator, and is associated with the voting power of one
// validator.
type Delegation struct {
	DelegatorAccount AccountID `json:"delegator_account" yaml:"delegator_account"`
	ValidatorAccount AccountID `json:"validator_account" yaml:"validator_account"`
	Shares           Dec       `json:"shares" yaml:"shares"`
}

// NewDelegation creates a new delegation object
func NewDelegation(delegatorAddr chainTypes.AccountID, validatorAddr chainTypes.AccountID, shares sdk.Dec) Delegation {
	return Delegation{
		DelegatorAccount: delegatorAddr,
		ValidatorAccount: validatorAddr,
		Shares:           shares,
	}
}

// MustMarshalDelegation returns the delegation bytes. Panics if fails
func MustMarshalDelegation(cdc *codec.Codec, delegation Delegation) []byte {
	return cdc.MustMarshalBinaryBare(&delegation)
}

// MustUnmarshalDelegation return the unmarshaled delegation from bytes.
// Panics if fails.
func MustUnmarshalDelegation(cdc *codec.Codec, value []byte) Delegation {
	delegation, err := UnmarshalDelegation(cdc, value)
	if err != nil {
		panic(err)
	}
	return delegation
}

// return the delegation
func UnmarshalDelegation(cdc *codec.Codec, value []byte) (delegation Delegation, err error) {
	err = cdc.UnmarshalBinaryBare(value, &delegation)
	return delegation, err
}

// nolint
func (d Delegation) Equal(d2 Delegation) bool {
	return bytes.Equal(d.DelegatorAccount.Bytes(), d2.DelegatorAccount.Bytes()) &&
		bytes.Equal(d.ValidatorAccount.Bytes(), d2.ValidatorAccount.Bytes()) &&
		d.Shares.Equal(d2.Shares)
}

// nolint - for Delegation
func (d Delegation) GetDelegatorAccountID() chainTypes.AccountID { return d.DelegatorAccount }
func (d Delegation) GetValidatorAccountID() chainTypes.AccountID { return d.ValidatorAccount }
func (d Delegation) GetDelegatorAddr() sdk.AccAddress {
	delegatorAccAddr, _ := d.DelegatorAccount.ToAccAddress()
	return delegatorAccAddr
}
func (d Delegation) GetValidatorAddr() sdk.ValAddress {
	validatorAccaddr, _ := d.ValidatorAccount.ToAccAddress()
	return sdk.ValAddress(validatorAccaddr)
}
func (d Delegation) GetShares() sdk.Dec { return d.Shares }

// String returns a human readable string representation of a Delegation.
func (d Delegation) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// Delegations is a collection of delegations
type Delegations []Delegation

func (d Delegations) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func NewUnbondingDelegationEntry(
	creationHeight int64,
	completionTime time.Time,
	balance sdk.Int) UnbondingDelegationEntry {
	return UnbondingDelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

// String implements the stringer interface for a UnbondingDelegationEntry.
func (e UnbondingDelegationEntry) String() string {
	out, _ := yaml.Marshal(e)
	return string(out)
}

// IsMature - is the current entry mature
func (e UnbondingDelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUnbondingDelegation - create a new unbonding delegation object
func NewUnbondingDelegation(
	delegatorAddr chainTypes.AccountID, validatorAddr chainTypes.AccountID,
	creationHeight int64, minTime time.Time, balance sdk.Int,
) UnbondingDelegation {

	return UnbondingDelegation{
		DelegatorAccount: delegatorAddr,
		ValidatorAccount: validatorAddr,
		Entries: []UnbondingDelegationEntry{
			NewUnbondingDelegationEntry(creationHeight, minTime, balance),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (ubd *UnbondingDelegation) AddEntry(creationHeight int64, minTime time.Time, balance sdk.Int) {
	entry := NewUnbondingDelegationEntry(creationHeight, minTime, balance)
	ubd.Entries = append(ubd.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (ubd *UnbondingDelegation) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// return the unbonding delegation
func MustMarshalUBD(cdc *codec.Codec, ubd UnbondingDelegation) []byte {
	return cdc.MustMarshalBinaryBare(&ubd)
}

// unmarshal a unbonding delegation from a store value
func MustUnmarshalUBD(cdc *codec.Codec, value []byte) UnbondingDelegation {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}
	return ubd
}

// unmarshal a unbonding delegation from a store value
func UnmarshalUBD(cdc *codec.Codec, value []byte) (ubd UnbondingDelegation, err error) {
	err = cdc.UnmarshalBinaryBare(value, &ubd)
	return ubd, err
}

// UnbondingDelegation stores all of a single delegator's unbonding bonds
// for a single validator in an time-ordered list
type UnbondingDelegation struct {
	DelegatorAccount AccountID                  `json:"delegator_account" yaml:"delegator_account"`
	ValidatorAccount AccountID                  `json:"validator_account" yaml:"validator_account"`
	Entries          []UnbondingDelegationEntry `json:"entries" yaml:"entries"`
}

func (ubd UnbondingDelegation) Equal(d2 UnbondingDelegation) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&ubd)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&d2)
	return bytes.Equal(bz1, bz2)
}

// String returns a human readable string representation of an UnbondingDelegation.
func (ubd UnbondingDelegation) String() string {
	out := fmt.Sprintf(`Unbonding Delegations between:
  Delegator:                 %s
  Validator:                 %s
	Entries:`, ubd.DelegatorAccount, ubd.ValidatorAccount)
	for i, entry := range ubd.Entries {
		out += fmt.Sprintf(`    Unbonding Delegation %d:
      Creation Height:           %v
      Min time to unbond (unix): %v
      Expected balance:          %s`, i, entry.CreationHeight,
			entry.CompletionTime, entry.Balance)
	}
	return out
}

// UnbondingDelegations is a collection of UnbondingDelegation
type UnbondingDelegations []UnbondingDelegation

func (ubds UnbondingDelegations) String() (out string) {
	for _, u := range ubds {
		out += u.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// UnbondingDelegationEntry defines an unbonding object with relevant metadata.
type UnbondingDelegationEntry struct {
	CreationHeight int64     `json:"creation_height,omitempty" yaml:"creation_height"`
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"`
	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"`
	Balance        sdk.Int   `json:"balance" yaml:"balance"`
}

// RedelegationEntry defines a redelegation object with relevant metadata.
type RedelegationEntry struct {
	CreationHeight int64     `json:"creation_height,omitempty" yaml:"creation_height"`
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"`
	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"`
	SharesDst      sdk.Dec   `json:"shares_dst" yaml:"shares_dst"`
}

func NewRedelegationEntry(
	creationHeight int64,
	completionTime time.Time,
	balance sdk.Int,
	sharesDst sdk.Dec) RedelegationEntry {
	return RedelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		SharesDst:      sharesDst,
	}
}

// String implements the Stringer interface for a RedelegationEntry object.
func (e RedelegationEntry) String() string {
	out, _ := yaml.Marshal(e)
	return string(out)
}

// IsMature - is the current entry mature
func (e RedelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// Redelegation contains the list of a particular delegator's redelegating bonds
// from a particular source validator to a particular destination validator.
type Redelegation struct {
	DelegatorAccount    AccountID           `json:"delegator_account" yaml:"delegator_account"`
	ValidatorSrcAccount AccountID           `json:"validator_src_account" yaml:"validator_src_account"`
	ValidatorDstAccount AccountID           `json:"validator_dst_account" yaml:"validator_dst_account"`
	Entries             []RedelegationEntry `json:"entries" yaml:"entries"`
}

func NewRedelegation(
	delegatorAddr chainTypes.AccountID, validatorSrcAddr, validatorDstAddr chainTypes.AccountID,
	creationHeight int64, minTime time.Time, balance sdk.Int, sharesDst sdk.Dec,
) Redelegation {

	return Redelegation{
		DelegatorAccount:    delegatorAddr,
		ValidatorSrcAccount: validatorSrcAddr,
		ValidatorDstAccount: validatorDstAddr,
		Entries: []RedelegationEntry{
			NewRedelegationEntry(creationHeight, minTime, balance, sharesDst),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (red *Redelegation) AddEntry(creationHeight int64, minTime time.Time, balance sdk.Int, sharesDst sdk.Dec) {
	entry := NewRedelegationEntry(creationHeight, minTime, balance, sharesDst)
	red.Entries = append(red.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (red *Redelegation) RemoveEntry(i int64) {
	red.Entries = append(red.Entries[:i], red.Entries[i+1:]...)
}

// MustMarshalRED returns the Redelegation bytes. Panics if fails.
func MustMarshalRED(cdc *codec.Codec, red Redelegation) []byte {
	return cdc.MustMarshalBinaryBare(&red)
}

// MustUnmarshalRED unmarshals a redelegation from a store value. Panics if fails.
func MustUnmarshalRED(cdc *codec.Codec, value []byte) Redelegation {
	red, err := UnmarshalRED(cdc, value)
	if err != nil {
		panic(err)
	}
	return red
}

// UnmarshalRED unmarshals a redelegation from a store value
func UnmarshalRED(cdc *codec.Codec, value []byte) (red Redelegation, err error) {
	err = cdc.UnmarshalBinaryBare(value, &red)
	return red, err
}

// nolint
// inefficient but only used in tests
func (red Redelegation) Equal(d2 Redelegation) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&red)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&d2)
	return bytes.Equal(bz1, bz2)
}

// String returns a human readable string representation of a Redelegation.
func (red Redelegation) String() string {
	out := fmt.Sprintf(`Redelegations between:
  Delegator:                 %s
  Source Validator:          %s
  Destination Validator:     %s
  Entries:
`,
		red.DelegatorAccount, red.ValidatorSrcAccount, red.ValidatorDstAccount,
	)

	for i, entry := range red.Entries {
		out += fmt.Sprintf(`    Redelegation Entry #%d:
      Creation height:           %v
      Min time to unbond (unix): %v
      Dest Shares:               %s
`,
			i, entry.CreationHeight, entry.CompletionTime, entry.SharesDst,
		)
	}

	return strings.TrimRight(out, "\n")
}

// Redelegations are a collection of Redelegation
type Redelegations []Redelegation

func (d Redelegations) String() (out string) {
	for _, red := range d {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ----------------------------------------------------------------------------
// Client Types

// DelegationResponse is equivalent to Delegation except that it contains a balance
// in addition to shares which is more suitable for client responses.
type DelegationResponse struct {
	Delegation
	Balance chainTypes.Coin `json:"balance" yaml:"balance"`
}

// NewDelegationResp creates a new DelegationResponse instance
func NewDelegationResp(
	delegatorAddr chainTypes.AccountID, validatorAddr chainTypes.AccountID, shares sdk.Dec, balance chainTypes.Coin,
) DelegationResponse {
	return DelegationResponse{
		Delegation: NewDelegation(delegatorAddr, validatorAddr, shares),
		Balance:    balance,
	}
}

// String implements the Stringer interface for DelegationResponse.
func (d DelegationResponse) String() string {
	return fmt.Sprintf("%s\n  Balance:   %s", d.Delegation.String(), d.Balance)
}

type delegationRespAlias DelegationResponse

// MarshalJSON implements the json.Marshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (d DelegationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal((delegationRespAlias)(d))
}

// UnmarshalJSON implements the json.Unmarshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (d *DelegationResponse) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, (*delegationRespAlias)(d))
}

// DelegationResponses is a collection of DelegationResp
type DelegationResponses []DelegationResponse

// String implements the Stringer interface for DelegationResponses.
func (d DelegationResponses) String() (out string) {
	for _, del := range d {
		out += del.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// RedelegationResponse is equivalent to a Redelegation except that its entries
// contain a balance in addition to shares which is more suitable for client
// responses.
type RedelegationResponse struct {
	Redelegation
	Entries []RedelegationEntryResponse `json:"entries" yaml:"entries"`
}

// NewRedelegationResponse crates a new RedelegationEntryResponse instance.
func NewRedelegationResponse(
	delegatorAddr, validatorSrc, validatorDst chainTypes.AccountID,
	entries []RedelegationEntryResponse) RedelegationResponse {
	return RedelegationResponse{
		Redelegation: Redelegation{
			DelegatorAccount:    delegatorAddr,
			ValidatorSrcAccount: validatorSrc,
			ValidatorDstAccount: validatorDst,
		},
		Entries: entries,
	}
}

// RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
// contains a balance in addition to shares which is more suitable for client
// responses.
type RedelegationEntryResponse struct {
	RedelegationEntry
	Balance sdk.Int `json:"balance"`
}

// NewRedelegationEntryResponse creates a new RedelegationEntryResponse instance.
func NewRedelegationEntryResponse(
	creationHeight int64, completionTime time.Time, sharesDst sdk.Dec,
	initialBalance, balance sdk.Int) RedelegationEntryResponse {
	return RedelegationEntryResponse{
		RedelegationEntry: NewRedelegationEntry(creationHeight, completionTime, initialBalance, sharesDst),
		Balance:           balance,
	}
}

// String implements the Stringer interface for RedelegationResp.
func (r RedelegationResponse) String() string {
	out := fmt.Sprintf(`Redelegations between:
  Delegator:                 %s
  Source Validator:          %s
  Destination Validator:     %s
  Entries:
`,
		r.DelegatorAccount, r.ValidatorSrcAccount, r.ValidatorDstAccount,
	)

	for i, entry := range r.Entries {
		out += fmt.Sprintf(`    Redelegation Entry #%d:
      Creation height:           %v
      Min time to unbond (unix): %v
      Initial Balance:           %s
      Shares:                    %s
      Balance:                   %s
`,
			i, entry.CreationHeight, entry.CompletionTime, entry.InitialBalance, entry.SharesDst, entry.Balance,
		)
	}

	return strings.TrimRight(out, "\n")
}

type redelegationRespAlias RedelegationResponse

// MarshalJSON implements the json.Marshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (r RedelegationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal((redelegationRespAlias)(r))
}

// UnmarshalJSON implements the json.Unmarshaler interface. This is so we can
// achieve a flattened structure while embedding other types.
func (r *RedelegationResponse) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, (*redelegationRespAlias)(r))
}

// RedelegationResponses are a collection of RedelegationResp
type RedelegationResponses []RedelegationResponse

func (r RedelegationResponses) String() (out string) {
	for _, red := range r {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
}
