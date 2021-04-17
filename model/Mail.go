package model

import (
	"encoding/json"
	"github.com/lenistwo/config"
	"github.com/lenistwo/util"
	"github.com/streadway/amqp"
	"strconv"
)

const (
	ContentType = "application/json"
)

type Mail struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Delay   int32  `json:"delay"`
}

func (mail *Mail) Publish(connection *amqp.Connection) {
	channel, err := connection.Channel()
	util.CheckError(err)
	defer channel.Close()

	marshaledMessage, err := json.Marshal(mail)
	util.CheckError(err)

	message := amqp.Publishing{
		ContentType: ContentType,
		Expiration:  strconv.FormatInt(int64(mail.Delay), 10),
		Body:        marshaledMessage,
	}

	err = channel.Publish(config.MailExchange, config.RoutingKey, config.Mandatory, config.Immediate, message)
	util.CheckError(err)
}
