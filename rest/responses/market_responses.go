package responses

import (
	"github.com/chaosknight/bitnet/entity"
)

type (
	Ticker struct {
		Basic
		Tickers []*entity.Ticker `json:"data,omitempty"`
	}
	IndexTicker struct {
		Basic
		IndexTickers []*entity.IndexTicker `json:"data,omitempty"`
	}
	OrderBook struct {
		Basic
		OrderBooks []*entity.OrderBook `json:"data,omitempty"`
	}
	Candle struct {
		Basic
		Candles []*entity.Candle `json:"data,omitempty"`
	}
	IndexCandle struct {
		Basic
		Candles []*entity.IndexCandle `json:"data,omitempty"`
	}
	Candleentity struct {
		Basic
		Candles []*entity.IndexCandle `json:"data,omitempty"`
	}
	Trade struct {
		Basic
		Trades []*entity.Trade `json:"data,omitempty"`
	}
	TotalVolume24H struct {
		Basic
		TotalVolume24Hs []*entity.TotalVolume24H `json:"data,omitempty"`
	}
	IndexComponent struct {
		Basic
		IndexComponents *entity.IndexComponent `json:"data,omitempty"`
	}
)
