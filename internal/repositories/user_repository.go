package repositories

import (
	"fmt"
	"payment-process/internal/database"
)

type User struct {
	Id int `db:"id"`
	Name string `db:"name"`
	LastName string `db:"last_name"`
	Email string `db:"email"`
	AccountId string `db:"account_id"`
}

type URepository struct {
	pers database.Persistence
}

func NewURepository(pers database.Persistence) *URepository {
	return &URepository{
		pers: pers,
	}
}


// This function `FindUserByAccountId` is a method of the `URepository` struct. It is used to find a
// user by their account ID in the database. Here is a breakdown of what the function does:
func (uRepository *URepository) FindUserByAccountId(accountId string) (*User, error) {
	db, err := uRepository.pers.DBConector.OpenConnect()

	if err != nil {
		return nil, fmt.Errorf("opening connection: %v", err)
	}

	defer db.Close()

	query := `SELECT * FROM users where account_id = ?`
	
	user := new(User)
	err = db.Get(user, query, accountId)

	if err != nil {
		return nil, fmt.Errorf("getting users: %v", err)
	}
	return user, nil
}