package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	queue, err := ch.QueueDeclare(
		queueName,
		queueType.Durable,
		!queueType.Durable,
		!queueType.Durable,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	err = ch.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		},
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return ch, queue, nil
}
