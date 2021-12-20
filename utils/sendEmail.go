package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(recipient string) {

  // Sender data.
  from := os.Getenv("EMAIL")
  password := os.Getenv("EMAIL_PASSWORD")

  // Receiver email address.
  to := []string{
    recipient,
  }

  // smtp server configuration.
  smtpHost := "smtp.gmail.com"
  smtpPort := "587"

  // Message.
  message := []byte("please sign up!!")
  
  // Authentication.
  auth := smtp.PlainAuth("", from, password, smtpHost)
  
  // Sending email.
  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println("Email Sent Successfully!")
}


