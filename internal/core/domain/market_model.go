package domain

import (
	"github.com/shopspring/decimal"
	mm "github.com/tdex-network/tdex-daemon/pkg/marketmaking"
	"github.com/tdex-network/tdex-daemon/pkg/marketmaking/formula"
)

// Market defines the Market entity data structure for holding an asset pair state
type Market struct {
	// AccountIndex links a market to a HD wallet account derivation.
	AccountIndex int
	BaseAsset    string
	QuoteAsset   string
	// Each Market has a different fee expressed in basis point of each swap
	Fee int64
	// if curretly open for trades
	Tradable bool
	// Market Making strategy
	Strategy mm.MakingStrategy
	// Pluggable Price of the asset pair.
	Price Prices
}

// OutpointWithAsset contains the transaction outpoint (tx hash and vout) along with the asset hash
type OutpointWithAsset struct {
	Asset string
	Txid  string
	Vout  int
}

// Prices ...
type Prices struct {
	// how much 1 base asset is valued in quote asset.
	BasePrice decimal.Decimal
	// how much 1 quote asset is valued in base asset
	QuotePrice decimal.Decimal
}

// StrategyType is the Market making strategy type
type StrategyType int32

// NewMarket returns an empty market with a reference to an account index.
// It is also mandatory to define a fee (in BP) for the market.
func NewMarket(positiveAccountIndex int, feeInBasisPoint int64) (*Market, error) {
	if err := validateAccountIndex(positiveAccountIndex); err != nil {
		return nil, err
	}

	if err := validateFee(feeInBasisPoint); err != nil {
		return nil, err
	}

	return &Market{
		AccountIndex: positiveAccountIndex,
		Fee:          feeInBasisPoint,
		Strategy:     mm.NewStrategyFromFormula(formula.BalancedReserves{}),
	}, nil
}
