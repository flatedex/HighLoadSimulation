package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

type User struct {
	Name string
	Password string
	Code int
}

func FailOnError(err error, msg string){
	if(err != nil) { 
		fmt.Println(msg)
		panic(err)
	}
}

func main() {	
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to connect to a channel")

	defer ch.Close()

	codeGeneratorQueue, err := ch.QueueDeclare(
		"codeGenerator",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue: codeGenerator")

	codeSenderQueue, err := ch.QueueDeclare(
		"codeSender",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue: codeSender")

	authQueue, err := ch.QueueDeclare(
		"auth",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue: auth")

	msg, err := ch.Consume(
		codeGeneratorQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer: codeGenerator")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var forever chan struct{}

	go func(){
		for d := range msg {
			fmt.Printf("Recieved a message: %s\n", d.Body)
			user := User{}
			err = json.Unmarshal([]byte(d.Body), &user)
			FailOnError(err, "Failed to convert JSON into object of user")

			user.Code = int(rand.Float64() * 100000)
			message := strconv.Itoa(user.Code)

			err = ch.PublishWithContext(
				ctx,
				"",
				codeSenderQueue.Name,
				false,
				false,
				amqp.Publishing{Body: []byte(message)},
			)
			FailOnError(err, "Failed to send message to codeSender")

			
			messageToAuth, err := json.Marshal(user)
			FailOnError(err, "Failed to create JSON from user")
			
			err = ch.PublishWithContext(
				ctx,
				"",
				authQueue.Name,
				false,
				false,
				amqp.Publishing{Body: []byte(messageToAuth)},
			)
			FailOnError(err, "Failed to send message to auth")
		}
	}()

	fmt.Println("Recieving messages, to exit press [Ctrl+C]")
	<- forever
}