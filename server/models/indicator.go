package models

import "gorm.io/gorm"

type Indicator struct {
	gorm.Model
	Name   string  // e.g) U.S Dollar Index
	Ticker string  `gorm:"uniqueIndex:uniq_idx_ticker"`     // e.g) "DXY"
	Value  float64 `gorm:"type:decimal(10,2);default:0.00"` // "100.21"
	Type   string  `gorm:"default:Fiat"`                    // "Fiat" | "Crypto" | "ETF"
	Domain string  // "www.tradingview.com"
}
