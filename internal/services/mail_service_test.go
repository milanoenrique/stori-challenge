package services_test

import (
	"payment-process/internal/repositories"
	"payment-process/internal/services"
	"testing"

	"github.com/go-gomail/gomail"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

type sMailServices struct {
	mock.Mock
}

type sSender struct {
	mock.Mock
}

func (s *sSender) Send(mail *gomail.Message) error {
	args:= s.Called(mail)
	return args.Error(0)
}

type sRepository struct {
	mock.Mock
}

// FindUserByAccountId implements services.IURepository.
func (s *sRepository) FindUserByAccountId(accountId string) (*repositories.User, error) {
	args := s.Called(accountId)
	return args.Get(0).(*repositories.User), args.Error(1)
}

// SendEmail implements services.IEmailSender.
func (s *sMailServices) SendEmail(accountNumber string, resumeTransaction services.ResumeAccount) error {
	args := s.Called(accountNumber, resumeTransaction)
	return args.Error(0)
}

func TestNewEmailSender(t *testing.T) {
	sMServices := new(sMailServices)
	sRepository := new(sRepository)
	sSender := new(sSender)

	user := repositories.User{}
	user.Id = 1
	user.AccountId = "123456"
	user.Name = "user"
	user.LastName = "test"
	user.Email = "test@email.com"

	sRepository.On("FindUserByAccountId", mock.Anything).Return(&repositories.User{
		Id:    1,
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}, nil)

	sSender.On("Send", mock.Anything).Return(nil)
	ms := services.NewEmailSender(sMServices, sRepository, sSender)

	resume := services.ResumeAccount{}
	resume.AverageCredit = 100.00
	resume.AverageDebit = -50.00
	resume.Total = 1000
	resume.TransactionsByMount = make(map[string]int)
	resume.TransactionsByMount["Jul"] = 10

	err := ms.SendEmail("123456.csv", resume)
	assert.NoError(t, err)
}
