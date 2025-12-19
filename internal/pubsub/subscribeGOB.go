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

// func encode(gl GameLog) ([]byte, error) {
//     var buf bytes.Buffer
//     enc := gob.NewEncoder(&buf)
//     err := enc.Encode(gl)
//     if err != nil {
//         return nil, err
//     }
//     return buf.Bytes(), nil
// }

// func decode(data []byte) (GameLog, error) {
// 	var log GameLog
//     buf := bytes.NewBuffer(data)
//     dec := gob.NewDecoder(buf)

//     // Decode the data from the buffer into the log variable
//     err := dec.Decode(&log)
//     if err != nil {
//         return GameLog{}, err
//     }

//	    return log, nil
//	}

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
	delCH, err := ch.Consume(queueName, "", false, false, false, false, nil)
	go func(delCH <-chan amqp.Delivery) {
		defer fmt.Print("> ")
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
