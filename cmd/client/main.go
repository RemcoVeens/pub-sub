package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const RMQConnectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril client...")

	amqpConn, err := amqp.Dial(RMQConnectionString)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer amqpConn.Close()
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Failed to get welcome message:", err)
		return
	}
	_, que, err := pubsub.DeclareAndBind(
		amqpConn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.SimpleQueueType{Durable: false},
	)
	fmt.Println(que.Name)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\rclose command received, shutting down now")
	os.Exit(0)
}
