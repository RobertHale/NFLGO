package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"spikenet.com/nflgo/internal/calculator"
	"spikenet.com/nflgo/internal/models"
)

const (
	KeyLocation = "../key.txt"
	APIURL      = "https://api.the-odds-api.com/v4/sports/americanfootball_nfl/odds?regions=us&markets=h2h&oddsFormat=american&apiKey="
)

func main() {
	var (
		keyByt    []byte
		keyVal    string
		resp      *http.Response
		oddsModel []models.SingleOdd
		client    = &http.Client{Timeout: 10 * time.Second}
		timeLimit = time.Now().Add(time.Hour * 24 * 7)

		err error
	)

	keyByt, err = os.ReadFile(KeyLocation)
	if err != nil {
		fmt.Printf("Error opening key file: %s\n", err.Error())
		os.Exit(1)
	}

	keyVal = strings.TrimSpace(string(keyByt))

	resp, err = client.Get(APIURL + keyVal)
	if err != nil {
		fmt.Printf("Error getting API: %s\n", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	// decode response body
	err = json.NewDecoder(resp.Body).Decode(&oddsModel)
	if err != nil {
		fmt.Printf("Error decoding JSON: %s\n", err.Error())
		os.Exit(1)
	}

	calc := calculator.New()
	for _, gOdd := range oddsModel {
		// skip if not this week
		if gOdd.CommenceTime.After(timeLimit) {
			continue
		}

		calc.AddGame(gOdd)
	}

	calc.PrintRanked()
}
