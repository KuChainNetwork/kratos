package rest

import (
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	rest "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// CommunityPoolSpendProposalReq defines a community pool spend proposal request body.
	CommunityPoolSpendProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title              string              `json:"title" yaml:"title"`
		Description        string              `json:"description" yaml:"description"`
		Recipient          chainType.AccountID `json:"recipient" yaml:"recipient"`
		Amount             sdk.Coins           `json:"amount" yaml:"amount"`
		Proposer           chainType.AccountID `json:"proposer" yaml:"proposer"`
		Deposit            sdk.Coins           `json:"deposit" yaml:"deposit"`
		ProposerAccAddress sdk.AccAddress      `json:"proposer_accaddress" yaml:"proposer_accaddress"`
	}
)
