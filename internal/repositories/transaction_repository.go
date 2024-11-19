package repositories

import (
	"context"
	"fmt"
	"payment-process/internal/database"
)

type Transaction struct {
	Id        int     `db:"id"`
	AccountId string  `db:"account_id"`
	Date      string  `db:"date"`
	Amount    float64 `db:"amount"`
}

type TXRepository struct {
	pers database.Persistence
}

func NewTXRepository(pers database.Persistence) *TXRepository {
	return &TXRepository{
		pers: pers,
	}
}

// This function `InsertTransaction` is responsible for inserting a transaction record into the
// database.
func (txRepository *TXRepository) InsertTransaction( tx *Transaction) error {
	db, err := txRepository.pers.DBConector.OpenConnect()

	if err != nil {
		return fmt.Errorf("opening connection: %v", err)
	}

	defer db.Close()

	query := `INSERT INTO transactions
	(id, account_id, date, amount) VALUES (?, ?, ?, ?)`

	_, err = db.ExecContext(context.Background(), query, tx.Id, tx.AccountId, tx.Date, tx.Amount)
	if err != nil {
		return fmt.Errorf("inserting transaction in database: %v", err)
	}
	return nil
}

func (txRepository *TXRepository) GetTransactions(accountId string) ([]Transaction, error){
	db, err := txRepository.pers.DBConector.OpenConnect()

	if err != nil {
		return nil, fmt.Errorf("opening connection: %v", err)
	}

	defer db.Close()
	query := `SELECT * FROM transactions where account_id = ?`

	txs := make([]Transaction, 0)

	 err = db.Select(&txs, query, accountId)

	 if err != nil {
		return nil, err
	 }

	 return txs, nil
}
