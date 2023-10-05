package calculator

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/shopspring/decimal"
	"spikenet.com/nflgo/internal/models"
)

type CalcArr []models.SingleGame

func (a CalcArr) Len() int           { return len(a) }
func (a CalcArr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CalcArr) Less(i, j int) bool { return a[i].GameTime.Before(a[j].GameTime) }

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
			GameTime:   odd.CommenceTime,
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

	*c = append(*c, game)
}

func (c *CalcArr) PrintRaw() {
	// sort array
	tempC := *c
	sort.Sort(tempC)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "HomeTeam\tHomeOdds\tAwayTeam\tAwayOdds\tGameTime\t")

	for _, game := range tempC {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", game.HomeTeam, game.HomePrice, game.AwayTeam, game.AwayPrice, game.GameTime)
	}

	w.Flush()
}
