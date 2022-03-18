package usecases

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mindaugasrmn/uphold-bot/pkg/domain"
	"github.com/mindaugasrmn/uphold-bot/pkg/helpers"
	repository "github.com/mindaugasrmn/uphold-bot/pkg/repo"
	"github.com/shopspring/decimal"
)

type Usecase interface {
	InitDB() error
	Ticker(q []domain.Query) error
}

type usecase struct {
	repo repository.TickerRepository
}

func NewUsecase(
	repo repository.TickerRepository,
) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) InitDB() error {
	return u.repo.Migrate()
}

func (u *usecase) Ticker(q []domain.Query) error {
	var wg sync.WaitGroup
	wg.Add(len(q))

	for _, v := range q {
		go func(v domain.Query, wg *sync.WaitGroup) {
			defer wg.Done()
			Tick(&v, u)
		}(v, &wg)
	}

	wg.Wait()
	return nil
}

func Tick(pair *domain.Query, u *usecase) {
	var priceChange domain.PriceOscillation
	priceChange.FirstTime = true
	finish := false

	// Fetch Upload Api each
	ticker := time.NewTicker(time.Second * time.Duration(pair.FetchInterval))
	defer ticker.Stop()
	for range ticker.C {
		errorIntent := 0
		data, code, err := helpers.HttpGET("https://api.uphold.com/v0/ticker/" + pair.CurrencyPair)
		if err != nil {
			// max 5 retries
			if errorIntent < 5 {
				errorIntent++
				continue
			}
			break
		}
		if *code != 200 {
			break
		}

		var res = &domain.Response{}
		helpers.DecodeResponseBody(data, res)
		finish, err = CheckPriceOscillation(u, *pair, *res, &priceChange)
		if err != nil {
			if errorIntent < 5 {
				errorIntent++
				continue
			}
			break
		}
		if finish {
			break
		}
	}
	return
}

func CheckPriceOscillation(u *usecase, q domain.Query, input domain.Response, obj *domain.PriceOscillation) (bool, error) {
	bid, err := decimal.NewFromString(input.Bid)
	if err != nil {
		return false, err
	}
	ask, err := decimal.NewFromString(input.Ask)
	if err != nil {
		return false, err
	}
	//check if price is valid
	if ask.IsZero() || ask.IsNegative() {
		return false, fmt.Errorf("invalid price received from api")
	}
	if bid.IsZero() || bid.IsNegative() {
		return false, fmt.Errorf("invalid price received from api")
	}

	//map Response  to PriceOscillation struct
	newData := domain.PriceOscillation{
		Ask:          ask,
		Bid:          bid,
		Currency:     input.Currency,
		CurrencyPair: q.CurrencyPair,
		Timestamp:    time.Now(),
	}

	if obj.FirstTime {
		obj.Ask = ask
		obj.Bid = bid
		obj.CurrencyPair = q.CurrencyPair
		obj.Currency = input.Currency
		obj.Timestamp = time.Now()
		obj.FirstTime = false
		return false, nil
	}

	//check if price changed more then fixed
	finish, err := CheckIfPriceChanged(u, &newData, obj, &q)
	if err != nil {
		return false, err
	}
	return finish, nil
}

func CheckIfPriceChanged(u *usecase, newData *domain.PriceOscillation, oldData *domain.PriceOscillation, q *domain.Query) (bool, error) {

	percentAsk := newData.Ask.Div(oldData.Ask).Sub(decimal.NewFromFloat(1)).RoundBank(4)

	percentBid := newData.Bid.Div(oldData.Bid).Sub(decimal.NewFromFloat(1)).RoundBank(4)

	msg := fmt.Sprint("PRICE OF ", newData.CurrencyPair)
	log.Println(newData.CurrencyPair, " ASK: ", newData.Ask, ", BID: ", newData.Bid, " | ", percentAsk, percentBid)

	// alert if price percentage changed
	if percentAsk.Abs().GreaterThan(q.PriceOsciliationInterval) {
		action := "INCRESED "
		if percentAsk.IsNegative() {
			action = "DECREASED "
		}
		msg = fmt.Sprint(msg, " ASK ", action, percentAsk.RoundBank(5).Abs().String(), "%")
		msg = msg + " | INIT: " + oldData.Ask.String() + ", " + oldData.Bid.String() + " | CUR: " + newData.Ask.String() + ", " + newData.Bid.String() + " | Oscillation: " + percentAsk.String() + ", " + percentBid.String()

		fmt.Println(msg)
		_, err := u.repo.Create(domain.PriceOscillationDB{
			InitAsk:             oldData.Ask,
			InitBid:             newData.Bid,
			CurAsk:              newData.Ask,
			CurBid:              newData.Bid,
			CurrencyPair:        newData.CurrencyPair,
			Timestamp:           newData.Timestamp,
			PriceOsciliationAsk: percentAsk,
			PriceOsciliationBid: percentBid,
			Type:                "ASK",
			Action:              action,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		return true, nil
	}

	if percentBid.Abs().GreaterThan(q.PriceOsciliationInterval) {
		action := " INCRESED "
		if percentBid.IsNegative() {
			action = " DECREASED "
		}
		msg = fmt.Sprint(msg, " BID ", action, percentBid.RoundBank(5).Abs().String(), "%")
		msg = msg + " | INIT: " + oldData.Ask.String() + ", " + oldData.Bid.String() + " | CUR: " + newData.Ask.String() + ", " + newData.Bid.String() + " | Oscillation: " + percentAsk.String() + ", " + percentBid.String()

		fmt.Println(msg)
		_, err := u.repo.Create(domain.PriceOscillationDB{
			InitAsk:             oldData.Ask,
			InitBid:             newData.Bid,
			CurAsk:              newData.Ask,
			CurBid:              newData.Bid,
			CurrencyPair:        newData.CurrencyPair,
			Timestamp:           newData.Timestamp,
			PriceOsciliationAsk: percentAsk,
			PriceOsciliationBid: percentBid,
			Type:                "BID",
			Action:              action,
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		return true, nil
	}
	return false, nil
}
