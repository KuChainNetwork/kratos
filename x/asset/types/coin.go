package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

type CoinDescription struct {
	Symbol      Name   `json:"symbol" yaml:"symbol"`                     // Symbol coin symbol name
	Creator     Name   `json:"creator" yaml:"creator"`                   // Creator coin creator account name
	Description []byte `json:"description,omitempty" yaml:"description"` // Description coin description info
}

// NewCoinDescription create coin description
func NewCoinDescription(creator, symbol Name, desc []byte) CoinDescription {
	return CoinDescription{
		Creator:     creator,
		Symbol:      symbol,
		Description: desc,
	}
}

func (m CoinDescription) String() string {
	res, _ := yaml.Marshal(m)
	return string(res)
}

// CoinStat state for a coin type
type CoinStat struct {
	Symbol        Name  `json:"symbol" yaml:"symbol"`   // Symbol coin symbol name
	Creator       Name  `json:"creator" yaml:"creator"` // Creator coin creator account name
	CreateHeight  int64 `json:"create_height,omitempty" yaml:"create_height"`
	Supply        Coin  `json:"supply" yaml:"supply"`         // Supply coin current supply
	MaxSupply     Coin  `json:"max_supply" yaml:"max_supply"` // MaxSupply coin max supply limit
	CanIssue      bool  `json:"can_issue,omitempty" yaml:"can_issue"`
	CanLock       bool  `json:"can_lock,omitempty" yaml:"can_lock"`
	CanBurn       bool  `json:"can_burn,omitempty" yaml:"can_burn"`
	IssueToHeight int64 `json:"issue_to_height,omitempty" yaml:"issue_to_height"`
	InitSupply    Coin  `json:"init_supply" yaml:"init_supply"` // InitSupply coin init supply, if issue_to_height is not zero, this will be the start supply for issue
}

// NewCoinStat creates a Coin status
func NewCoinStat(ctx sdk.Context, creator, symbol Name, maxSupply Coin) CoinStat {
	return CoinStat{
		Creator:      creator,
		Symbol:       symbol,
		CreateHeight: ctx.BlockHeight(),
		Supply:       NewCoin(maxSupply.Denom, types.NewInt(0)),
		MaxSupply:    maxSupply,
	}
}

// SetOpt set coin optional
func (c *CoinStat) SetOpt(canIssue, canLock, canBurn bool, issue2Height int64, initSupply Coin) error {
	if err := CheckCoinStatOpts(c.CreateHeight, canIssue, canLock, issue2Height, initSupply, c.MaxSupply); err != nil {
		return err
	}

	c.CanIssue = canIssue
	c.CanLock = canLock
	c.CanBurn = canBurn
	c.IssueToHeight = issue2Height
	c.InitSupply = initSupply
	return nil
}

func CheckCoinStatOpts(createHeight int64, canIssue, canLock bool, issue2Height int64, init, max Coin) error {
	if !canIssue {
		if (issue2Height != 0) || (!init.IsZero()) {
			return ErrAssetCoinMustCanIssueWhenIssueByBlock
		}
	}

	if issue2Height != 0 {
		if !max.IsGTE(init) {
			return ErrAssetCoinMustSupplyNeedGTInitSupply
		}

		if issue2Height <= createHeight {
			return ErrAssetIssueToHeightMustGTCurrentHeight
		}
	}

	if !(init.IsZero() || max.IsGTE(init)) {
		return ErrAssetMaxSupplyShouldGTEInitSupply
	}

	return nil
}

func (c *CoinStat) GetCurrentMaxSupplyLimit(currentHeight int64) types.Coin {
	if c.IssueToHeight == 0 {
		return c.MaxSupply
	}

	denom := c.MaxSupply.Denom

	if (currentHeight <= c.CreateHeight) || (c.IssueToHeight <= c.CreateHeight) {
		return types.NewCoin(denom, sdk.ZeroInt())
	}

	if currentHeight >= c.IssueToHeight {
		return c.MaxSupply
	}

	issuedHeight := currentHeight - c.CreateHeight
	allIssueHeight := c.IssueToHeight - c.CreateHeight

	if !c.MaxSupply.IsGTE(c.InitSupply) {
		return c.MaxSupply
	}

	needIssue := c.MaxSupply.Sub(c.InitSupply)
	addedIssue := needIssue.Amount.MulRaw(issuedHeight).QuoRaw(allIssueHeight)

	return c.InitSupply.Add(types.NewCoin(needIssue.Denom, addedIssue))
}

func (m CoinStat) String() string {
	res, _ := yaml.Marshal(m)
	return string(res)
}
