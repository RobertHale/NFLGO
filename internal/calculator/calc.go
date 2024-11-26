package calculator

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"spikenet.com/nflgo/internal/models"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/slices"
)

type CalcArr []models.SingleGame

func (c CalcArr) Len() int           { return len(c) }
func (c CalcArr) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CalcArr) Less(i, j int) bool { return c[i].GameTime.Before(c[j].GameTime) }

func New() CalcArr {
	return CalcArr{}
}

func (c *CalcArr) AddGame(odd models.SingleOdd) {
	var (
		game = models.SingleGame{
			HomeTeam:   odd.HomeTeam,
			AwayTeam:   odd.AwayTeam,
			HomePrice:  decimal.NewFromInt(0),
			AwayPrice:  decimal.NewFromInt(0),
			PriceCount: decimal.NewFromInt(0),
			GameTime:   odd.CommenceTime.Local(),
		}
		inc = decimal.NewFromInt(1)
	)

	for _, bm := range odd.Bookmakers {
		// odds are not there or unexpected
		if len(bm.Markets) == 0 || len(bm.Markets[0].Outcomes) < 2 {
			continue
		}

		oc1 := bm.Markets[0].Outcomes[0]
		oc2 := bm.Markets[0].Outcomes[1]

		switch game.HomeTeam {
		case oc1.Name:
			game.HomePrice = game.HomePrice.Add(decimal.NewFromInt(int64(oc1.Price)))
			game.AwayPrice = game.AwayPrice.Add(decimal.NewFromInt(int64(oc2.Price)))
		case oc2.Name:
			game.HomePrice = game.HomePrice.Add(decimal.NewFromInt(int64(oc2.Price)))
			game.AwayPrice = game.AwayPrice.Add(decimal.NewFromInt(int64(oc1.Price)))
		}

		game.PriceCount = game.PriceCount.Add(inc)
	}

	game.HomePrice = game.HomePrice.Div(game.PriceCount)
	game.AwayPrice = game.AwayPrice.Div(game.PriceCount)
	game.PriceDiff = game.HomePrice.Sub(game.AwayPrice).Abs()

	*c = append(*c, game)
}

func (c *CalcArr) PrintRaw() {
	// sort array
	tempC := *c
	sort.Sort(tempC)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "HomeTeam\tHomeOdds\tAwayTeam\tAwayOdds\tRank\tGameTime\t")

	for _, game := range tempC {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t\n", game.HomeTeam, game.HomePrice.Round(1), game.AwayTeam, game.AwayPrice.Round(1), game.Rank, game.GameTime.Format(time.RFC1123))
	}

	w.Flush()
}

func (c *CalcArr) PrintRanked() {
	var (
		rankMap = make(map[float64]models.SingleGame)
		diffArr []float64
		tempC   = *c
	)

	for _, game := range tempC {
		diffF := game.PriceDiff.InexactFloat64()
		diffArr = append(diffArr, diffF)
		rankMap[diffF] = game
	}

	slices.Sort(diffArr)

	tempC = CalcArr{}
	for i, diff := range diffArr {
		temp := rankMap[diff]
		temp.Rank = i + 1
		tempC = append(tempC, temp)
	}

	tempC.PrintRaw()
}
