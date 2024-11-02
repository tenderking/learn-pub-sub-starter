package pubsub

import (
	"bytes"
	"encoding/gob"
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

func jsonUnmarshaller[T any](data []byte) (T, error) {
	var target T
	err := json.Unmarshal(data, &target)
	return target, err
}

func gobDecoder[T any](data []byte) (T, error) {
	var target T
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&target)
	return target, err
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange string,
	queueName string,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
) error {

	ch, q, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		int(simpleQueueType),
	)
	failOnError(err, "failed to declare and bind")
	ch.Qos(10, 0, false)

	consumer := ""
	deliveries, err := ch.Consume(q.Name, consumer, false, false, false, false, nil)

	failOnError(err, "failed to register a consumer")

	go func() {
		defer ch.Close()
		for msg := range deliveries {
			target, err := unmarshaller(msg.Body)

			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			ack := handler(target)
			switch ack {
			case Ack:
				msg.Ack(false)
				fmt.Println("Acked")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println("Nacked and requeued")
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("Nacked and discarded")
			}

		}
	}()
	return nil

}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange string,
	queueName string,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		jsonUnmarshaller,
	)
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange string,
	queueName string,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		gobDecoder[T],
	)
}
