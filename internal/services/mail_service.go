package services

import (
	"fmt"
	"os"
	"payment-process/internal/repositories"

	"strings"

	"github.com/go-gomail/gomail"
)

type ISender interface {
    Send(mail *gomail.Message) error
}

type RealSender struct {
	iSender ISender
}

type IEmailSender interface {
	SendEmail(accountNumber string, resumeTransaction  ResumeAccount) error
}

type IURepository interface {
	FindUserByAccountId(accountId string) (*repositories.User, error)
}

type EmailSender struct {
	IEmailSender IEmailSender
	repo IURepository
	iSender ISender
}

type Email struct {
	To string
	Body []string
}

func NewEmailSender(iEmailSender IEmailSender, repo IURepository, iSender ISender) EmailSender{
	return EmailSender {
		IEmailSender: iEmailSender,
		repo: repo,
		iSender: iSender,
	}
}

const _subject = "Account details"


func (e *EmailSender)SendEmail(accountNumber string, resumeTransaction ResumeAccount) error {
	accountId := strings.Split(accountNumber,".")[0]
	user, err  := e.repo.FindUserByAccountId(accountId)

	if err!= nil {
		return err
	}

	email := new(Email)
	email.To = user.Email

	transactionByMonth := `"<tr><td>Number of transactions in <month>: </td><td><numberOfTransaction></td></tr>"`
	

	for month, txByMonth := range resumeTransaction.TransactionsByMount {
		transactionByMonth = strings.Replace(transactionByMonth, "<month>", month, -1)
		transactionByMonth = strings.Replace(transactionByMonth, "<numberOfTransaction>",  fmt.Sprintf("%v", txByMonth), -1)
		transactionByMonth = transactionByMonth + transactionByMonth
	}

	emailFrom := os.Getenv("emailFrom")

	mail := gomail.NewMessage()
	mail.SetHeader("From", emailFrom)
	mail.SetHeader("To", email.To)
	mail.SetHeader("Subject", _subject)
	mail.Embed("./images/image.png", gomail.Rename("logo.png"), gomail.SetHeader(map[string][]string{
		"Content-ID": {"<logoImage>"},
	}))

	mail.SetHeader("Subject", "Story Resume Account")
	body :=  `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Story Resume Account</title>
	</head>
	<body>
		<center><img src="cid:logoImage" alt="Logo" style="width:200px;height:auto;"></center>
		<h1 style="color: blue;">Your Resume Account</h1>
		<p>Hi <username> heres is your resume account.</p>
		<table>
			<tr><td>Total Balance is: </td><td><totoalBalance></td><tr>
			<tr><td>Average Debit Amount </td><td><averageDebit></td><tr>
			<tr><td>Average Credit Amount </td><td><averageCrebit></td><tr>
			<transactionsByMonth>
		</table>


	</body>
	</html>
	`
	body = strings.Replace(body, "<transactionsByMonth>", transactionByMonth, -1)
	mail.SetBody("text/html",body)

	err =e.iSender.Send(mail)

	if err != nil {
		return err
	}

	return nil
}

func  (rs *RealSender)Send(mail *gomail.Message) error {
	emailFrom := os.Getenv("emailFrom")
	pass := os.Getenv("password")
	d := gomail.NewDialer("smtp.gmail.com", 587, emailFrom, pass)

	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
