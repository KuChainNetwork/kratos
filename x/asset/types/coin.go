package types

import (
	"gopkg.in/yaml.v2"

	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// NewCoinStat creates a Coin status
func NewCoinStat(ctx sdk.Context, creator, symbol Name, maxSupply Coin) CoinStat {
	return CoinStat{
		Creator:      creator,
		Symbol:       symbol,
		CreateHeight: ctx.BlockHeight(),
		Supply:       NewCoin(maxSupply.GetDenom(), sdk.NewInt(0)),
		MaxSupply:    maxSupply,
	}
}

// SetOpt set coin optional
func (c *CoinStat) SetOpt(canIssue, canLock bool, issue2Height int64, initSupply Coin) error {
	if err := CheckCoinStatOpts(c.CreateHeight, canIssue, canLock, issue2Height, initSupply, c.MaxSupply); err != nil {
		return err
	}

	c.CanIssue = canIssue
	c.CanLock = canLock
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
