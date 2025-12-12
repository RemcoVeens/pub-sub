package main

import (
	"fmt"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const RMQConnectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	gamelogic.PrintServerHelp()
	amqpConn, err := amqp.Dial(RMQConnectionString)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer amqpConn.Close()
	channel, err := amqpConn.Channel()
	if err != nil {
		fmt.Println("Failed to open channel:", err)
		return
	}

	for true {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "quit":
			os.Exit(0)
		case "pause":
			fmt.Println("pausing")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Println("Failed to publish message:", err)
				return
			}
		case "resume":
			fmt.Println("resuming")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		default:
			fmt.Println("Invalid command")
			gamelogic.PrintServerHelp()
		}
	}

}
