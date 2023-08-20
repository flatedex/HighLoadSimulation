package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/json"
	"math/rand"
	"strconv"
)

const PORT = 8000

type User struct {
	Name string
	Password string
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
	FailOnError(err, "Failed to declare a queue")

	codeSenderQueue, err := ch.QueueDeclare(
		"codeSender",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	msg, err := ch.Consume(
		codeGeneratorQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func(){
		for d := range msg {
			fmt.Printf("Recieved a message: %s\n", d.Body)
			user := User{}
			err = json.Unmarshal([]byte(d.Body), &user)
			FailOnError(err, "Failed to convert object into JSON")

			code := int(rand.Float64() * 100000)
			message := strconv.Itoa(code)

			err = ch.Publish(
				"",
				codeSenderQueue.Name,
				false,
				false,
				amqp.Publishing{Body: []byte(message)},
			)
			FailOnError(err, "Failed to send message to codeSender")
		}
	}()

	fmt.Println("Recieving messages, to exit press [Ctrl+C]")
	<- forever
}