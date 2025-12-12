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

func HandlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	defer fmt.Print("> ")
	return func(ps routing.PlayingState) {
		gs.HandlePause(ps)
	}
}

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
	_, _, err = pubsub.DeclareAndBind(
		amqpConn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.SimpleQueueType{Durable: false},
	)
	if err != nil {
		fmt.Println("Failed to declare and bind:", err)
		return
	}

	gameState := gamelogic.NewGameState(username)
	err = pubsub.SubscribeJSON(amqpConn, string(routing.ExchangePerilDirect), fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey, pubsub.SimpleQueueType{Durable: false}, HandlerPause(gameState))
	if err != nil {
		fmt.Println("Failed to subscribe:", err)
		return
	}

	for true {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err = gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println("Failed to spawn:", err)
			}
		case "move":
			move, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println("Failed to move:", err)
				continue
			}
			fmt.Println("Moved to:", move)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			os.Exit(0)
		default:
			fmt.Println("Unknown command:", words[0])
			continue
		}
	}

	// fmt.Println(que.Name)
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// fmt.Println("\rclose command received, shutting down now")
	// os.Exit(0)
}
