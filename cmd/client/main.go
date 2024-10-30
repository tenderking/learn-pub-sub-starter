package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall" // Import syscall for SIGTERM

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tenderking/learn-pub-sub-starter/internal/gamelogic"
	"github.com/tenderking/learn-pub-sub-starter/internal/pubsub"
	"github.com/tenderking/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")
	addr := "amqp://guest:guest@localhost:5672/"
	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	conn, err := amqp.Dial(addr)
	if err != nil {
		fmt.Println("Error connecting to RabbitMQ")
		return
	}
	defer conn.Close()
	rabbitmqChannel, err := conn.Channel()
	if err != nil {
		fmt.Println("Error creating RabbitMQ channel")
		return
	}
	defer rabbitmqChannel.Close()
	username, _ := gamelogic.ClientWelcome()
	queueName := []string{routing.PauseKey, username}
	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, strings.Join(queueName, "."), routing.PauseKey, 0)
	// Notify for both SIGINT and SIGTERM
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
