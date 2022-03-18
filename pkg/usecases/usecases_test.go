package usecases

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/mindaugasrmn/uphold-bot/pkg/domain"
	repo "github.com/mindaugasrmn/uphold-bot/pkg/repo"
	"github.com/shopspring/decimal"
)

var tickerUse Usecase
var tickerRepo repo.TickerRepository

const fileName = "sqlite_test.db"

type RequestTable struct {
	TestName string `json:"test_name"`
	Filter   []domain.Query
}

func setupTickerUseCase() {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	tickerRepo = repo.NewTickerRepository(db)
	tickerUse = NewUsecase(tickerRepo)
	tickerUse.InitDB()
}

func TestPriceOscillationBot(t *testing.T) {
	setupTickerUseCase()

	validCases, err := SingleCurrencyAlert()
	assertNoErr(t, err)
	RunTestsCase(t, validCases)

	validCases, err = MultipleCurrencyAlert()
	assertNoErr(t, err)
	RunTestsCase(t, validCases)

	invalidCases, err := InvalidSingleCurrencyAlert()
	assertNoErr(t, err)
	RunTestsInvalidCase(t, invalidCases)
	err = os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}

}

func SingleCurrencyAlert() ([]RequestTable, error) {

	priceOscilationInterval1 := decimal.NewFromFloat(0.0001)

	pair1 := domain.Query{
		CurrencyPair:             "BTC-EUR",
		FetchInterval:            2,
		PriceOsciliationInterval: priceOscilationInterval1,
	}

	q := []domain.Query{}
	q = append(q, pair1)

	validCases := []RequestTable{
		{TestName: "Test multiple currency alert", Filter: q},
	}
	return validCases, nil
}

func MultipleCurrencyAlert() ([]RequestTable, error) {
	priceOscilationInterval1 := decimal.NewFromFloat(0.0003)
	priceOscilationInterval2 := decimal.NewFromFloat(0.0004)
	priceOscilationInterval3 := decimal.NewFromFloat(0.0005)

	pair1 := domain.Query{
		CurrencyPair:             "BTC-EUR",
		FetchInterval:            7,
		PriceOsciliationInterval: priceOscilationInterval1,
	}

	pair2 := domain.Query{
		CurrencyPair:             "ETH-EUR",
		FetchInterval:            3,
		PriceOsciliationInterval: priceOscilationInterval2,
	}
	pair3 := domain.Query{
		CurrencyPair:             "ADA-EUR",
		FetchInterval:            2,
		PriceOsciliationInterval: priceOscilationInterval3,
	}
	q := []domain.Query{}
	q = append(q, pair1, pair2, pair3)

	validCases := []RequestTable{
		{TestName: "", Filter: q},
	}
	return validCases, nil
}

func RunTestsCase(t *testing.T, validCases []RequestTable) {
	for _, x := range validCases {
		t.Run(x.TestName, func(t *testing.T) {
			err := tickerUse.Ticker(x.Filter)
			if err != nil {
				t.Errorf("Error failed with error: %v", err)
			}
		})
	}
}

func RunTestsInvalidCase(t *testing.T, invalidCases []RequestTable) {
	for _, x := range invalidCases {
		t.Run(x.TestName, func(t *testing.T) {

			err := tickerUse.Ticker(x.Filter)
			if err == nil {
				t.Errorf("Should return error and returns %v", err)
			}
		})
	}
}

func InvalidSingleCurrencyAlert() ([]RequestTable, error) {
	priceOscilationInterval1 := decimal.NewFromFloat(0.0004)

	pair1 := domain.Query{
		CurrencyPair:             "BTC-RUB",
		FetchInterval:            3,
		PriceOsciliationInterval: priceOscilationInterval1,
	}

	q := []domain.Query{}
	q = append(q, pair1)

	validCases := []RequestTable{
		{TestName: "Test invalid pair", Filter: q},
	}
	return validCases, nil
}

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}
