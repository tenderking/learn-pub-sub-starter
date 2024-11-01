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
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Error creating channel")
		return
	}
	defer ch.Close()

	username, _ := gamelogic.ClientWelcome()
	queueName := []string{routing.PauseKey, username}
	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		strings.Join(queueName, "."),
		routing.PauseKey,
		pubsub.Durable,
		handlerPause(gameState),
	)
	if err != nil {
		fmt.Println("Error subscribing to queue", err)
		return
	}
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		strings.Join([]string{routing.ArmyMovesPrefix, gameState.Player.Username}, "."),
		routing.ArmyMovesPrefix+".*",
		pubsub.Transient,
		handlerMoves(gameState, ch),
	)
	if err != nil {
		fmt.Println("Error subscribing to queue", err)
		return
	}
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.Durable,
		handlerWar(gameState, ch),
	)
	if err != nil {
		fmt.Println("Error subscribing to queue", err)
		return
	}

	for {
		word := gamelogic.GetInput()
		if word == nil {
			continue
		}

		if word[0] == "spawn" {
			fmt.Println("Spawning a new player...")
			err := gameState.CommandSpawn(word)
			if err != nil {
				fmt.Println("Error spawning unit:", err)
			}
		}
		if word[0] == "move" {
			fmt.Println("Moving a player...")
			_, err := gameState.CommandMove(word)
			if err != nil {
				fmt.Println("Error moving unit:", err)

			}
			var units []gamelogic.Unit
			for _, unit := range gameState.GetPlayerSnap().Units {
				units = append(units, unit)
			}
			err = pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+gameState.GetPlayerSnap().Username,
				gamelogic.ArmyMove{
					Player:     gameState.GetPlayerSnap(),
					Units:      units,
					ToLocation: gameState.GetPlayerSnap().Units[0].Location,
				},
			)
			if err != nil {
				fmt.Println("Error publishing message")
			}
			fmt.Println("Player moved!")

		}
		if word[0] == "status" {
			fmt.Println("Checking the status of the game...")
			gameState.CommandStatus()
		}
		if word[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		}
		if word[0] == "help" {
			gamelogic.PrintClientHelp()
		}

		if word[0] == "quit" {

			gamelogic.PrintQuit()
			break
		}

		if word[0] != "spawn" && word[0] != "move" && word[0] != "status" && word[0] != "spam" && word[0] != "quit" && word[0] != "help" {
			fmt.Println("Invalid command. Please try again.")
			continue
		}
	}

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
