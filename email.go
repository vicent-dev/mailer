package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
	"strconv"
)
type Body struct {
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Recaptcha string `json:"recaptcha"`
}

func (b Body) isValid() bool {
	return b.Email != "" && b.Subject != "" && b.Body != "" && b.Recaptcha != ""
}

func (b Body) EmailMessage() string {
	return "Email: " + b.Email + "<br>Subject: " + b.Subject + "<br>Message: " + b.Body
}
func sendEmail(to string, subject string, body string) bool {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"),
		port, os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		log.Print(err)
		return false
	}
	return true
}

func sendEmailHandler(c *gin.Context) {
	var body Body
	json.NewDecoder(c.Request.Body).Decode(&body)
	if !body.isValid() {
		writErrorResponse(c.Writer, errors.New("Wrong params"), http.StatusBadRequest)
		return
	}
	if ok, recaptchaErrors := recaptcha(body.Recaptcha); !ok {
		writErrorResponse(c.Writer, errors.New("Invalid recaptcha"), http.StatusBadRequest)
		log.Println(recaptchaErrors)
		return
	}
	success := sendEmail(os.Getenv("DEFAULT_SEND"), "Email presupuesto de " + body.Email, body.EmailMessage())

	if !success {
		writErrorResponse(c.Writer, errors.New("Can't send email"), 0)
		return
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return
}


