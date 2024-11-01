package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {

	ch, q, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		int(simpleQueueType),
	)
	failOnError(err, "failed to declare and bind")

	consumer := ""
	deliveries, err := ch.Consume(q.Name, consumer, false, false, false, false, nil)

	failOnError(err, "failed to register a consumer")
	unmarshaller := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}

	go func() {
		defer ch.Close()
		for msg := range deliveries {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			ack := handler(target)
			if ack == Ack {
				msg.Ack(false)
				fmt.Println("Acked")
			}
			if ack == NackRequeue {
				msg.Nack(false, true)
				fmt.Println("Nacked and requeued")

			}
			if ack == NackDiscard {
				msg.Nack(false, false)
				fmt.Println("Nacked and discarded")
			}
		}
	}()
	return nil

}
