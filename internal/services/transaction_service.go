package services

import (
	
	"payment-process/internal/repositories"
	"payment-process/pkg/utils"
	"strconv"
	"strings"
)

const _dirPath = "../../transactions" // Ruta del directorio con los archivos CSV
type Month int
const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	Dececember
)

var monthMap = map[Month]string{
	January:      "January",
	February:    "February",
	March:      "March",
	April:      "April",
	May:       "May",
	June:      "June",
	July:      "July",
	August:     "August",
	September: "September",
	October:    "October",
	November:  "November",
	Dececember:  "December",
}

type ITXRepository interface {
	InsertTransaction(tx *repositories.Transaction) error
	GetTransactions(accountId string) ([]repositories.Transaction, error)
}

type TransactionService struct {
	repo ITXRepository
}

type ResumeAccount struct {
	Total float64
	TransactionsByMount map[string]int
	AverageDebit float64
	AverageCredit float64
}

func NewTransactionService(tRepository ITXRepository) TransactionService{
	return TransactionService{
		repo: tRepository,
	}
}

func (ts *TransactionService) SaveTransaction(file string) error {
	
	records, err := utils.ReadCsv(_dirPath+"/"+file)
	if err != nil {
		return err
	}

	for _, record := range records[1:]{
		id, err := strconv.Atoi(record[0])
		
		if err != nil {
			return err
		}

		date := record[1]

		amount, err := strconv.ParseFloat(record[2], 64)

		if err != nil {
			return err
		}

		accountId := strings.Split(file, ".")[0]
		tx:= new(repositories.Transaction)
		tx.Id = id
		tx.Amount = amount
		tx.Date = date
		tx.AccountId = accountId
		
		err = ts.repo.InsertTransaction(tx)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (ts *TransactionService) GetResumeTransactions(fileName string) (*ResumeAccount, error) {
	accountId := strings.Split(fileName, ".")[0]
	txs,err  := ts.repo.GetTransactions(accountId)
	if err != nil {
		return nil, err
	}

	resumeAccount := new(ResumeAccount)
	credit := 0.00
	debit := 0.00
	for _, tx := range txs {
		resumeAccount.Total += tx.Amount
		if tx.Amount > 0 {
			resumeAccount.AverageCredit  += tx.Amount
			credit++
		} else if tx.Amount < 0 {
			resumeAccount.AverageDebit += tx.Amount
			debit++
		}

		monthString := strings.Split(tx.Date, "/")[0]

		monthNumber, err:= strconv.Atoi(monthString) 

		if err != nil {
			return nil, err
		}

		month := monthMap[Month(monthNumber)]

		if resumeAccount.TransactionsByMount == nil {
			resumeAccount.TransactionsByMount = make(map[string]int)
		}

		resumeAccount.TransactionsByMount[month]++	
	}
	resumeAccount.AverageCredit = resumeAccount.AverageCredit/credit
	resumeAccount.AverageDebit = resumeAccount.AverageDebit/debit
	
	return resumeAccount, nil
}

