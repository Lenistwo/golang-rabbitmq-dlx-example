package config

import (
	"github.com/lenistwo/util"
	"github.com/streadway/amqp"
)

const (
	DelayedMailQueueName   = "q.delayed.mail"
	MailQueueName          = "q.mail"
	DeadLetterExchangeName = "ex.delayed.mail"
	MailExchange           = "ex.mail"
	ExchangeType           = "direct"
	RoutingKey             = ""
	Durable                = false
	Exclusive              = false
	DeleteWhenUnused       = false
	NoWait                 = false
	Mandatory              = false
	Immediate              = false
)

func CreateRabbitStructure(amqConnection *amqp.Connection) {
	declareMailQueue(amqConnection)
	declareDeadLetterExchange(amqConnection)
	declareDelayedMailQueue(amqConnection)
	declareMailExchange(amqConnection)
}

func declareDelayedMailQueue(amqConnection *amqp.Connection) {
	channel, err := amqConnection.Channel()
	util.CheckError(err)
	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = DeadLetterExchangeName
	_, err = channel.QueueDeclare(DelayedMailQueueName, Durable, Exclusive, DeleteWhenUnused, NoWait, args)
	util.CheckError(err)
	defer channel.Close()
}

func declareMailQueue(amqConnection *amqp.Connection) {
	channel, err := amqConnection.Channel()
	util.CheckError(err)
	_, err = channel.QueueDeclare(MailQueueName, Durable, Exclusive, DeleteWhenUnused, NoWait, nil)
	util.CheckError(err)
	defer channel.Close()
}

func declareDeadLetterExchange(amqConnection *amqp.Connection) {
	channel, err := amqConnection.Channel()
	util.CheckError(err)
	err = channel.ExchangeDeclare(DeadLetterExchangeName, ExchangeType, Durable, Exclusive, DeleteWhenUnused, NoWait, nil)
	err = channel.QueueBind(MailQueueName, RoutingKey, DeadLetterExchangeName, NoWait, nil)
	util.CheckError(err)
	defer channel.Close()
}

func declareMailExchange(amqConnection *amqp.Connection) {
	channel, err := amqConnection.Channel()
	util.CheckError(err)
	err = channel.ExchangeDeclare(MailExchange, ExchangeType, Durable, Exclusive, DeleteWhenUnused, NoWait, nil)
	err = channel.QueueBind(DelayedMailQueueName, RoutingKey, MailExchange, NoWait, nil)
	util.CheckError(err)
	defer channel.Close()
}
