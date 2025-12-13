package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) (err error) {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	delCH, err := ch.Consume(queueName, "", false, false, false, false, nil)
	go func(delCH <-chan amqp.Delivery) {
		for mess := range delCH {
			var val T
			err := json.Unmarshal(mess.Body, &val)
			if err != nil {
				fmt.Println("Failed to unmarshal:", err)
				return
			}
			ackType := handler(val)
			switch ackType {
			case Ack:
				fmt.Println("Acknowledged")
				mess.Ack(false)
			case NackRequeue:
				fmt.Println("Nacked and Requeued")
				mess.Nack(false, true)
			case NackDiscard:
				fmt.Println("Nacked and Discarded")
				mess.Nack(false, false)
			default:
				fmt.Printf("Warning: Unknown AckType %d returned. Defaulting to Ack.\n", ackType)
				mess.Ack(false)
			}
		}
	}(delCH)
	return
}
