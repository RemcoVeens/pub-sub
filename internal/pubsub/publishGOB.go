package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGOB[t any](ch *amqp.Channel, exchange, key string, val t) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false, false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        buf.Bytes(),
		})
}

func PublishGamelog(ch *amqp.Channel, username, mess string) error {
	return PublishGOB(
		ch,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.GameLogSlug, username),
		mess,
	)
}
