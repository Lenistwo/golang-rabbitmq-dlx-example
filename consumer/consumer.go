package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lenistwo/config"
	"github.com/lenistwo/model"
	"github.com/lenistwo/util"
	"github.com/streadway/amqp"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

const (
	AutoAck  = true
	NoLocal  = false
	Consumer = ""
)

var (
	amqConnection *amqp.Connection
)

func init() {
	err := godotenv.Load(".env")
	util.CheckError(err)
	amqConnection = establishRabbitConnection()
	config.CreateRabbitStructure(amqConnection)
}

func main() {
	fmt.Printf("Started Consumer For Queue %s\n", config.MailQueueName)
	defer amqConnection.Close()
	channel, err := amqConnection.Channel()
	util.CheckError(err)

	consume, err := channel.Consume(config.MailQueueName, Consumer, AutoAck, config.Exclusive, NoLocal, config.NoWait, nil)
	util.CheckError(err)

	forever := make(chan bool)

	go func() {
		for message := range consume {
			var mail model.Mail
			err = json.Unmarshal(message.Body, &mail)
			util.CheckError(err)
			sendMail(mail)
			fmt.Println("Mail Sent!")
		}
	}()

	<-forever
}

func sendMail(mail model.Mail) {
	message := gomail.NewMessage()
	message.SetHeader("To", mail.To)
	message.SetHeader("From", mail.From)
	message.SetHeader("Subject", mail.Title)
	message.SetBody("text/plain", mail.Content)

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(message)
	util.CheckError(err)
}

func establishRabbitConnection() *amqp.Connection {
	dial, err := amqp.Dial(prepareAMQPConnectionURL())
	util.CheckError(err)
	return dial
}

func prepareAMQPConnectionURL() string {
	return "amqp://" + os.Getenv("RABBIT_USERNAME") + ":" +
		os.Getenv("RABBIT_PASSWORD") + "@" + os.Getenv("RABBIT_HOST") + ":" + os.Getenv("RABBIT_PORT") + "/"
}
