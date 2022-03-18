package repository

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/mindaugasrmn/uphold-bot/pkg/domain"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type tickerRepository struct {
	db *sql.DB
}

// NewAMLRepository initialization of repo
func NewTickerRepository(
	db *sql.DB,
) TickerRepository {
	return &tickerRepository{
		db: db,
	}
}

type TickerRepository interface {
	Migrate() error
	Create(website domain.PriceOscillationDB) (*domain.PriceOscillationDB, error)
	All() ([]domain.PriceOscillationDB, error)
}

func (r *tickerRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS data (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        init_ask TEXT NOT NULL UNIQUE,
        init_bid TEXT NOT NULL,
        cur_ask TEXT NOT NULL,
		cur_bid TEXT NOT NULL,
		price_osciliation_ask TEXT NOT NULL,
		price_osciliation_bid TEXT NOT NULL,
		currency_pair TEXT NOT NULL,
		timestamp DATE NOT NULL,
		type TEXT NOT NULL,
		action TEXT NOT NULL
    );
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *tickerRepository) Create(data domain.PriceOscillationDB) (*domain.PriceOscillationDB, error) {
	res, err := r.db.Exec("INSERT INTO data(init_ask, init_bid, cur_ask,cur_bid,price_osciliation_ask,price_osciliation_bid,currency_pair, timestamp,type,action ) values(?,?,?,?,?,?,?,?,?,?)",
		data.InitAsk, data.InitBid, data.CurAsk, data.CurBid, data.PriceOsciliationAsk, data.PriceOsciliationBid, data.CurrencyPair, data.Timestamp, data.Type, data.Action)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	data.ID = id

	return &data, nil
}

func (r *tickerRepository) All() ([]domain.PriceOscillationDB, error) {
	rows, err := r.db.Query("SELECT * FROM data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []domain.PriceOscillationDB
	for rows.Next() {
		var data domain.PriceOscillationDB
		if err := rows.Scan(&data.InitAsk, &data.InitBid, &data.CurAsk, &data.CurBid, &data.PriceOsciliationAsk, &data.PriceOsciliationBid, &data.CurrencyPair, &data.Timestamp, &data.Type, &data.Action); err != nil {
			return nil, err
		}
		all = append(all, data)
	}
	return all, nil
}
