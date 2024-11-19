package main

import (
	"fmt"
	"os"
	"payment-process/internal/config"
	"payment-process/internal/database"
	"payment-process/internal/repositories"
	"payment-process/internal/services"
	"payment-process/pkg/utils/logs"
	"strings"
	"sync"
)

const (
	_logPath = "./logs/"
	DirPath = "./transactions" // Ruta del directorio con los archivos CSV
)

func main(){

	fmt.Println("Job Start")
	
	start()
}

func start(){
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logLevel := logs.SetLoggerLevel(logs.ErrorLevel)
	logDir := "/Users/enrique.milano/Desktop/source/go/payment-process/cmd/job" //TODO arreglar
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}
	logger, err := logs.InitializeLog(logDir, logLevel)

	if err != nil {
		panic(err)
	}

	dbCnf := new(dbConfig)
	dbCnf.dbName = conf.Database.Database
	dbCnf.dbPassword = os.Getenv("DB_PASSWORD")

	dbCnf.dbPort = conf.Database.Port
	dbCnf.dbServer = conf.Database.Host
	dbCnf.dbUsername = conf.Database.User



	conn, err := initDBConnection(dbCnf)

	if err != nil {
		panic(err)
	}

	persistence := database.NewPersistence(&conn)

	tRepository := repositories.NewTXRepository(persistence)
	uRepository := repositories.NewURepository(persistence)

	txService := services.NewTransactionService(tRepository)

	eSender := new(services.EmailSender)

	realSend := new(services.RealSender)

	emailService := services.NewEmailSender(eSender, uRepository, realSend)

	numWorkers := 5
	

	accountNumbers, err := os.ReadDir(DirPath)

	if err != nil {
		logger.Fatal(fmt.Errorf("reading folder with csv: %v", err))
	}

	fileChan := make(chan string, len(accountNumbers))
	var wg sync.WaitGroup

	for _, accountNumber := range accountNumbers {
		if !accountNumber.IsDir() {
			if strings.Contains(accountNumber.Name(), "csv"){
				fileChan <- accountNumber.Name()
			} 
		}
	}
	close(fileChan)


	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for fileName := range fileChan {
				err = txService.SaveTransaction(fileName)
				if err != nil {
					logger.Error(fmt.Errorf("saving transaction: %v", err))
				}
				resumeTransaction, err  := txService.GetResumeTransactions(fileName)
				if err != nil {
					logger.Error(fmt.Errorf("getting resume transactions: %v", err))
				}

				err = emailService.SendEmail(fileName, *resumeTransaction)
				if err != nil {
					logger.Error(fmt.Errorf("sending email: %v", err))
				}
				os.Remove(fileName)
			}
		}(i)
	}

	wg.Wait()
	
	logger.Info("proccess complete")
}
type dbConfig struct {
	dbName     string
	dbPort     string
	dbServer   string
	dbUsername string
	dbPassword string
}

func initDBConnection(conf *dbConfig) (database.Connector, error) {
	dbManage, err := database.NewConnectionManager(conf.dbName, conf.dbPort, conf.dbServer, conf.dbUsername, conf.dbPassword)
	if err != nil {
		return database.Connector{}, fmt.Errorf("initializing database connection: %v", err)
	}

	db, err := dbManage.OpenConnect()
	if err != nil {
		return database.Connector{}, fmt.Errorf("conecting database: %v pass: %v", err, conf.dbPassword)
	}

	err = db.Ping()
	if err != nil {
		return database.Connector{}, fmt.Errorf("db ping: %v", err)
	}
	return dbManage, nil
}

