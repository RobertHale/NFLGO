package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// /////////////////////////
// JSON Response From API
// /////////////////////////
type SingleOdd struct {
	ID           string       `json:"id"`
	SportKey     string       `json:"sport_key"`
	SportTitle   string       `json:"sport_title"`
	CommenceTime time.Time    `json:"commence_time"`
	HomeTeam     string       `json:"home_team"`
	AwayTeam     string       `json:"away_team"`
	Bookmakers   []Bookmakers `json:"bookmakers"`
}
type Outcomes struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
type Markets struct {
	Key        string     `json:"key"`
	LastUpdate time.Time  `json:"last_update"`
	Outcomes   []Outcomes `json:"outcomes"`
}
type Bookmakers struct {
	Key        string    `json:"key"`
	Title      string    `json:"title"`
	LastUpdate time.Time `json:"last_update"`
	Markets    []Markets `json:"markets"`
}

// /////////////////////////////
// Internal calculation holder
// /////////////////////////////
type SingleGame struct {
	HomeTeam   string
	AwayTeam   string
	HomePrice  decimal.Decimal
	AwayPrice  decimal.Decimal
	PriceCount decimal.Decimal
	PriceDiff  decimal.Decimal
	Rank       int
	GameTime   time.Time
}
