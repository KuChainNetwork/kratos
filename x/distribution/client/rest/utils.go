package rest

import (
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// CommunityPoolSpendProposalReq defines a community pool spend proposal request body.
	CommunityPoolSpendProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title              string         `json:"title" yaml:"title"`
		Description        string         `json:"description" yaml:"description"`
		Recipient          AccountID      `json:"recipient" yaml:"recipient"`
		Amount             Coins          `json:"amount" yaml:"amount"`
		Proposer           AccountID      `json:"proposer" yaml:"proposer"`
		Deposit            Coins          `json:"deposit" yaml:"deposit"`
		ProposerAccAddress sdk.AccAddress `json:"proposer_accaddress" yaml:"proposer_accaddress"`
	}
)
