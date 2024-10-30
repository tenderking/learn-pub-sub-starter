package main

import (
	"fmt"
	"os"
	"os/signal"
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

	gamelogic.PrintServerHelp()

	rabbitmqChannel, err := conn.Channel()
	if err != nil {
		fmt.Println("Error creating RabbitMQ channel")
		return
	}

	for {
		word := gamelogic.GetInput()
		if word == nil {
			continue
		}
		if word[0] == "pause" {
			fmt.Println("Pausing the game...")
			err = pubsub.PublishJSON(
				rabbitmqChannel,
				routing.ExchangePerilDirect,
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
				rabbitmqChannel,
				routing.ExchangePerilDirect,
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

	defer conn.Close()
	defer rabbitmqChannel.Close()

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
