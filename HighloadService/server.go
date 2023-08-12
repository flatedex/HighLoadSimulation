package main

import (
	"fmt"
	"net/http"
	amqp "github.com/rabbitmq/amqp091-go"
)

const PORT = 8000

type Message struct {
	Name string
	Message string
}

func FailOnError(err error, msg string){
	if(err != nil) { 
		fmt.Println(msg)
		panic(err)
	}
}

func GetDataFromClient(writer http.ResponseWriter, request *http.Request){
	defer writer.Write([]byte(http.StatusText(200)))

	err := request.ParseForm()
	if(err != nil) { panic(err) }

	var message Message

	message.Name = request.FormValue("Name")
	message.Message = request.FormValue("Message")
}

func main() {	
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to connect to a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"myQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	msg, err := ch.Consume(
		q.Name,
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
		}
	}()

	fmt.Println("Recieving messages, to exit press [Ctrl+C]")
	<- forever
}