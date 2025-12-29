package main

import (
	"log"
	"time"

	"booking-app/internal/models"

	mail "github.com/xhit/go-simple-mail"
)

// listenForMail listen the mail data from the channel
func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

// sendMsg will send email
func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	server.Username = "test@example.com"
	server.Password = "examplepass"

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent successfully...")
	}
}
