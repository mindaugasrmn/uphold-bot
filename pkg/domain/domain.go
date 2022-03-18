package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Query struct {
	CurrencyPair             string          `json:"currennncy_pair"`
	FetchInterval            int             `json:"fetch_interval"`
	PriceOsciliationInterval decimal.Decimal `json:"price_osciliation_interval"`
}

type PriceOscillation struct {
	Ask                 decimal.Decimal `json:"ask"`
	Bid                 decimal.Decimal `json:"bid"`
	PriceOscilationAsk  decimal.Decimal `json:"price_osciliation_ask"`
	PriceOsciliationBid decimal.Decimal `json:"price_osciliation_bid"`
	Currency            string          `json:"currency"`
	CurrencyPair        string          `json:"currency_pair"`
	Timestamp           time.Time       `json:"timestamp"`
	FirstTime           bool            `json:"first_time"`
}

type PriceOscillationDB struct {
	ID                  int64           `json:"id"`
	InitAsk             decimal.Decimal `json:"init_ask"`
	InitBid             decimal.Decimal `json:"init_bid"`
	CurAsk              decimal.Decimal `json:"cur_ask"`
	CurBid              decimal.Decimal `json:"cur_bid"`
	PriceOsciliationAsk decimal.Decimal `json:"price_osciliation_ask"`
	PriceOsciliationBid decimal.Decimal `json:"price_osciliation_bid"`
	CurrencyPair        string          `json:"currency_pair"`
	Timestamp           time.Time       `json:"timestamp"`
	Type                string          `json:"type"`
	Action              string          `json:"action"`
}

type Response struct {
	Ask      string `json:"ask"`
	Bid      string `json:"bid"`
	Currency string `json:"currency"`
}
