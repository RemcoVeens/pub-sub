package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

const RMQConnectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")
	amqpConn, err := amqp.Dial(RMQConnectionString)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer amqpConn.Close()

	fmt.Println("Connected to RabbitMQ")
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\rclose command received, shutting down now")
	os.Exit(0)
}
