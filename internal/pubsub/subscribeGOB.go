package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGOB[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	unmarshaller func([]byte) (T, error),
) (err error) {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	ch.Qos(10, 0, false)
	delCH, err := ch.Consume(queueName, "", false, false, false, false, nil)
	go func(delCH <-chan amqp.Delivery) {
		for mess := range delCH {
			var val string
			dec := gob.NewDecoder(bytes.NewBuffer(mess.Body))
			err := dec.Decode(&val)

			if err != nil {
				fmt.Println("Failed to unmarshal:", err)
				return
			}
			gamelogic.WriteLog(routing.GameLog{
				Username:    "",
				Message:     val,
				CurrentTime: time.Now(),
			})
			mess.Ack(false)
		}
	}(delCH)
	return
}
