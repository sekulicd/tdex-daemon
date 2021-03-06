package application

import (
	"github.com/shopspring/decimal"
	"github.com/tdex-network/tdex-daemon/internal/core/domain"
)

// SwapInfo is the data struct returned by ListSwap RPC.
type SwapInfo struct {
	Status           int32
	AmountP          uint64
	AssetP           string
	AmountR          uint64
	AssetR           string
	MarketFee        Fee
	RequestTimeUnix  uint64
	AcceptTimeUnix   uint64
	CompleteTimeUnix uint64
	ExpiryTimeUnix   uint64
}

// MarketInfo is the data struct returned by ListMarket RPC.
type MarketInfo struct {
	Market       Market
	Fee          Fee
	Tradable     bool
	StrategyType int
}

type Market struct {
	BaseAsset  string
	QuoteAsset string
}

type MarketWithFee struct {
	Market
	Fee
}

// Fee is the market fee percentage in basis point:
// 	- 0,01% -> 1 bp
//	- 1,00% -> 100 bp
//	- 99,99% -> 9999 bp
type Fee struct {
	BasisPoint int64
}

type MarketWithPrice struct {
	Market
	Price
}

type Price struct {
	BasePrice  decimal.Decimal
	QuotePrice decimal.Decimal
}

type PriceWithFee struct {
	Price
	Fee
	Amount uint64
	Asset  string
}

type MarketStrategy struct {
	Market
	Strategy domain.StrategyType
}

type Balance struct {
	BaseAmount  int64
	QuoteAmount int64
}

type BalanceWithFee struct {
	Balance
	Fee
}

type BalanceInfo struct {
	TotalBalance       uint64
	ConfirmedBalance   uint64
	UnconfirmedBalance uint64
}

type WithdrawMarketReq struct {
	Market
	BalanceToWithdraw Balance
	MillisatPerByte   int64
	Address           string
	Push              bool
}

type ReportMarketFee struct {
	CollectedFees              []FeeInfo
	TotalCollectedFeesPerAsset map[string]int64
}

type AddressAndBlindingKey struct {
	Address     string
	BlindingKey string
}

type FeeInfo struct {
	TradeID     string
	BasisPoint  int64
	Asset       string
	Amount      uint64
	MarketPrice decimal.Decimal
}

type TxOutpoint struct {
	Hash  string
	Index int
}
