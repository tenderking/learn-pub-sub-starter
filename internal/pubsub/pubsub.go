package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error{

      body, err := json.Marshal(val)
      if err != nil {
          return err
      }
      err = ch.ExchangeDeclare(
    "logs",   // name
    "fanout", // type
    true,     // durable
    false,    // auto-deleted
    false,    // internal
    false,    // no-wait
    nil,      // arguments
  )
  failOnError(err, "Failed to declare an exchange")

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  err = ch.PublishWithContext(ctx,
    "logs", // exchange
    "",     // routing key
    false,  // mandatory
    false,  // immediate
    amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
    })
    return err
}
