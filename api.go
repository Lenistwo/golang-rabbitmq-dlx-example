package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lenistwo/config"
	"github.com/lenistwo/model"
	"github.com/lenistwo/util"
	"github.com/streadway/amqp"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	connection *amqp.Connection
)

func CreateRouter(amqConnection *amqp.Connection) *gin.Engine {
	connection = amqConnection
	server := gin.Default()
	server.POST("/send", sendMail)
	server.GET("/queue", queueStats)
	return server
}

func queueStats(c *gin.Context) {
	request, err := http.NewRequest("GET", buildQueueURL(), nil)
	util.CheckError(err)
	request.Header.Set("Authorization", prepareAuthorizationHeader())
	client := &http.Client{}
	apiResponse, err := client.Do(request)

	if err != nil {
		c.JSON(500, gin.H{"Message": "Internal Server Error"})
		return
	}

	body, err := ioutil.ReadAll(apiResponse.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)

	if err != nil {
		c.JSON(500, gin.H{"Message": "Internal Server Error"})
		return
	}

	c.JSON(200, response)
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

func buildQueueURL() string {
	return "http://" + os.Getenv("RABBIT_HOST") + ":" + os.Getenv("RABBIT_MANAGEMENT_PORT") + "/api/queues/%2f/" + config.DelayedMailQueueName
}

func prepareAuthorizationHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("RABBIT_USERNAME")+":"+os.Getenv("RABBIT_PASSWORD")))
}
