package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall" // Import syscall for SIGTERM

	// Import the pubsub package

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tenderking/learn-pub-sub-starter/internal/gamelogic"
	"github.com/tenderking/learn-pub-sub-starter/internal/pubsub"
	"github.com/tenderking/learn-pub-sub-starter/internal/routing"
)

func main() {
	addr := "amqp://guest:guest@localhost:5672/"
	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	conn, err := amqp.Dial(addr)
	if err != nil {
		fmt.Println("Error connecting to RabbitMQ")
		return
	}
	defer conn.Close()
	publishCh, err := conn.Channel()
	if err != nil {
		fmt.Println("Error creating channel")
		return
	}

	gamelogic.PrintServerHelp()

	key := []string{routing.GameLogSlug, "*"}

	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, strings.Join(key, "."), int(pubsub.Durable))
	if err != nil {
		fmt.Println("Error declaring and binding queue", err)
		return
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	for {
		word := gamelogic.GetInput()
		if word == nil {
			continue
		}
		if word[0] == "pause" {
			fmt.Println("Pausing the game...")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				fmt.Println("Error publishing message")
			}

		}
		if word[0] == "resume" {
			fmt.Println("Resuming the game...")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				fmt.Println("Error publishing message")
			}
		}
		if word[0] == "quit" {
			break
		}
	}
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done // Wait for the signal handler to signal completion
	fmt.Println("exiting")
}
