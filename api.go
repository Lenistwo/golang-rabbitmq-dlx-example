package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lenistwo/model"
	"github.com/lenistwo/util"
	"github.com/streadway/amqp"
	"io"
)

var (
	connection *amqp.Connection
)

func CreateRouter(amqConnection *amqp.Connection) *gin.Engine {
	connection = amqConnection
	server := gin.Default()
	server.POST("/send", sendMail)
	return server
}

func sendMail(context *gin.Context) {
	body := context.Request.Body
	defer body.Close()

	object, err := io.ReadAll(body)
	util.CheckError(err)

	var mail model.Mail
	err = json.Unmarshal(object, &mail)

	if err != nil {
		context.JSON(422, gin.H{"Message": "Unprocessable Entity"})
		return
	}

	mail.Publish(connection)
	context.JSON(200, gin.H{"Message": "Scheduled Mail"})
}
