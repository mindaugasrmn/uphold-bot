package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mindaugasrmn/uphold-bot/pkg/domain"
	"github.com/mindaugasrmn/uphold-bot/pkg/helpers"
	repository "github.com/mindaugasrmn/uphold-bot/pkg/repo"
	usecases "github.com/mindaugasrmn/uphold-bot/pkg/usecases"
	"github.com/shopspring/decimal"
)

const fileName = "sqlite.db"

func main() {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewTickerRepository(db)
	use := usecases.NewUsecase(repo)

	err = use.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	q := []domain.Query{}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf(`Insert a currency pair! {string: 'BTC-EUR', 'ETH-EUR' etc}: `)
		scanner.Scan()
		fmt.Printf("\n")
		currencyPair := strings.ToUpper(scanner.Text())

		fmt.Printf(`Insert a fetch interval in seconds! (integer: 5, 10 etc): `)
		scanner.Scan()
		fmt.Printf("\n")

		interval, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf(`Invalid Input. Insert a Fetch Interval in Seconds! (integer, like: 5, 10 etc): `)
			scanner.Scan()
			fmt.Printf("\n")
			interval, err = strconv.Atoi(scanner.Text())
			log.Fatal(err)
		}

		fmt.Printf(`Insert a price osciliation Interval! (float: 0.1, 0.02 etc): `)
		scanner.Scan()
		fmt.Printf("\n")
		priceOscilationInterval, err := decimal.NewFromString(scanner.Text())
		if err != nil {
			fmt.Printf("Invalid input. Insert a price osciliation interval! (float: 0.1, 0.03 etc): ")
			scanner.Scan()
			priceOscilationInterval, err = decimal.NewFromString(scanner.Text())
			log.Fatal(err)
		}

		pair := domain.Query{
			CurrencyPair:             currencyPair,
			FetchInterval:            interval,
			PriceOsciliationInterval: priceOscilationInterval,
		}

		_, code, err := helpers.HttpGET("https://api.uphold.com/v0/ticker/" + pair.CurrencyPair)
		if err == nil {
			if *code == 200 {
				q = append(q, pair)
			} else {
				fmt.Printf(`Can't add currency pair as it is not met in uphold!\n\n`)
			}
		}

		fmt.Printf(`Add another currency pair? (y/n): `)
		scanner.Scan()
		if strings.ToUpper(scanner.Text()) == "N" {
			break
		} else if strings.ToUpper(scanner.Text()) != "Y" {
			fmt.Printf(`Invalid Input. Add another currency pair? (y/n): `)
			scanner.Scan()
			fmt.Printf("\n\n")
		}

	}
	fmt.Printf("\n\nStarting ticker...\n")
	err = use.Ticker(q)
	if err != nil {
		log.Fatal("Error")
	}
	fmt.Println("All completed, exiting")

}
