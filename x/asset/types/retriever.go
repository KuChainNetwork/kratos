package types

import (
	"fmt"
)

// NodeQuerier is an interface that is satisfied by types that provide the QueryWithData method
type NodeQuerier interface {
	// QueryWithData performs a query to a Tendermint node with the provided path
	// and a data payload. It returns the result and height of the query upon success
	// or an error if the query fails.
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

// AssetRetriever defines the properties of a type that can be used to
// retrieve accounts.
type AssetRetriever struct {
	querier NodeQuerier
}

// NewAssetRetriever init a new AssetRetriever instance.
func NewAssetRetriever(querier NodeQuerier) AssetRetriever {
	return AssetRetriever{querier: querier}
}

// GetCoin queries for coin for a account
func (ar AssetRetriever) GetCoin(acc AccountID, creator, symbol Name) (Coin, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryCoinParams(acc, creator, symbol))
	if err != nil {
		return Coin{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryCoin), bs)
	if err != nil {
		return Coin{}, height, err
	}

	var coinData Coin
	if err := ModuleCdc.UnmarshalJSON(res, &coinData); err != nil {
		return Coin{}, height, err
	}

	return coinData, height, nil
}

// GetCoins queries for coins for a account
func (ar AssetRetriever) GetCoins(acc AccountID) (Coins, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryCoinsParams(acc))
	if err != nil {
		return Coins{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryCoins), bs)
	if err != nil {
		return Coins{}, height, err
	}

	var coinData Coins
	if err := ModuleCdc.UnmarshalJSON(res, &coinData); err != nil {
		return Coins{}, height, err
	}

	return coinData, height, nil
}

// GetCoin queries for coin for a account
func (ar AssetRetriever) GetCoinPower(acc AccountID, creator, symbol Name) (Coin, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryCoinPowerParams(acc, creator, symbol))
	if err != nil {
		return Coin{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryCoinPower), bs)
	if err != nil {
		return Coin{}, height, err
	}

	var coinData Coin
	if err := ModuleCdc.UnmarshalJSON(res, &coinData); err != nil {
		return Coin{}, height, err
	}

	return coinData, height, nil
}

// GetCoinPowers queries for coins powers for a account
func (ar AssetRetriever) GetCoinPowers(acc AccountID) (Coins, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryCoinPowersParams(acc))
	if err != nil {
		return Coins{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryCoinPowers), bs)
	if err != nil {
		return Coins{}, height, err
	}

	var coinData Coins
	if err := ModuleCdc.UnmarshalJSON(res, &coinData); err != nil {
		return Coins{}, height, err
	}

	return coinData, height, nil
}

// GetLockedCoins queries for coins locked for a account
func (ar AssetRetriever) GetLockedCoins(acc AccountID) (QueryLockedCoinsResponse, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryLockedCoinsParams(acc))
	if err != nil {
		return QueryLockedCoinsResponse{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryCoinLocked), bs)
	if err != nil {
		return QueryLockedCoinsResponse{}, height, err
	}

	var coinData QueryLockedCoinsResponse
	if err := ModuleCdc.UnmarshalJSON(res, &coinData); err != nil {
		return QueryLockedCoinsResponse{}, height, err
	}

	return coinData, height, nil
}

type GetCoinStatResponse struct {
	CoinStat

	CoinDescription       string `json:"description"`
	CurrentMaxSupplyLimit Coin   `json:"current_max_supply"`
}

func (ar AssetRetriever) GetCoinStat(creator, symbol Name) (GetCoinStatResponse, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryCoinStatParams(creator, symbol))
	if err != nil {
		return GetCoinStatResponse{}, 0, err
	}

	res, height, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryCoinStat), bs)
	if err != nil {
		return GetCoinStatResponse{}, height, err
	}

	resDesc, _, err := ar.querier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryCoinDescription), bs)
	if err != nil {
		return GetCoinStatResponse{}, height, err
	}

	var coinDescData CoinDescription
	if err := ModuleCdc.UnmarshalJSON(resDesc, &coinDescData); err != nil {
		return GetCoinStatResponse{}, height, err
	}

	var coinData CoinStat
	if err := ModuleCdc.UnmarshalJSON(res, &coinData); err != nil {
		return GetCoinStatResponse{}, height, err
	}

	resData := GetCoinStatResponse{
		CoinStat:              coinData,
		CoinDescription:       string(coinDescData.Description),
		CurrentMaxSupplyLimit: coinData.GetCurrentMaxSupplyLimit(height),
	}

	return resData, height, nil
}
