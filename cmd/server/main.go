package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall" // Import syscall for SIGTERM

	// Import the pubsub package

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tenderking/learn-pub-sub-starter/internal/pubsub"
)

func main() {
    addr := "amqp://guest:guest@localhost:5672/"
    signalChan  := make(chan os.Signal, 1)
    done := make(chan bool, 1)
    conn, err := amqp.Dial(addr)
    if err != nil {
        fmt.Println("Error connecting to RabbitMQ")
        return
    }
		rabbitmqChannel, err := conn.Channel()
		if err != nil {
			fmt.Println("Error creating RabbitMQ channel")
			return
		}

     err = pubsub.PublishJSON(rabbitmqChannel, "logs", "", "Hello, World!") 
    if err != nil {
        fmt.Println("Error publishing message:", err)
        return 
    }
    fmt.Println("Starting Peril server...")

    // Notify for both SIGINT and SIGTERM
    signal.Notify(signalChan , syscall.SIGINT, syscall.SIGTERM) 

    defer conn.Close()
		defer rabbitmqChannel.Close()

    go func() {
        sig := <-signalChan 
        fmt.Println()
        fmt.Println(sig)
        done <- true
    }()

    fmt.Println("awaiting signal")
    <-done  // Wait for the signal handler to signal completion
    fmt.Println("exiting")
}