package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorSigningInfo defines the signing info for a validator
type ValidatorSigningInfo struct {
	Address sdk.ConsAddress `json:"address,omitempty"`
	// height at which validator was first a candidate OR was unjailed
	StartHeight int64 `json:"start_height,omitempty" yaml:"start_height"`
	// index offset into signed block bit array
	IndexOffset int64 `json:"index_offset,omitempty" yaml:"index_offset"`
	// timestamp validator cannot be unjailed until
	JailedUntil time.Time `json:"jailed_until" yaml:"jailed_until"`
	// whether or not a validator has been tombstoned (killed out of validator set)
	Tombstoned bool `json:"tombstoned,omitempty"`
	// missed blocks counter (to avoid scanning the array every time)
	MissedBlocksCounter int64 `json:"missed_blocks_counter,omitempty" yaml:"missed_blocks_counter"`
}

// NewValidatorSigningInfo creates a new ValidatorSigningInfo instance
func NewValidatorSigningInfo(
	condAddr sdk.ConsAddress, startHeight, indexOffset int64,
	jailedUntil time.Time,
	tombstoned bool, missedBlocksCounter int64) ValidatorSigningInfo {
	return ValidatorSigningInfo{
		Address:             condAddr,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
	}
}

// String implements the stringer interface for ValidatorSigningInfo
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Index Offset:          %d
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d`,
		i.Address, i.StartHeight, i.IndexOffset, i.JailedUntil,
		i.Tombstoned, i.MissedBlocksCounter)
}

// unmarshal a validator signing info from a store value
func UnmarshalValSigningInfo(cdc *codec.Codec, value []byte) (signingInfo ValidatorSigningInfo, err error) {
	err = cdc.UnmarshalBinaryBare(value, &signingInfo)
	return signingInfo, err
}
