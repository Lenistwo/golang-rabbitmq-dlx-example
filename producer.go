package main

import (
	"github.com/joho/godotenv"
	"github.com/lenistwo/config"
	"github.com/lenistwo/util"
	"github.com/streadway/amqp"
	"os"
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
	defer amqConnection.Close()
	err := CreateRouter(amqConnection).Run(":" + os.Getenv("ROUTER_PORT"))
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
