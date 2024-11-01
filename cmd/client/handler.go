package main

import (
	"fmt"
	"time"

	// Import syscall for SIGTERM

	"github.com/rabbitmq/amqp091-go"
	"github.com/tenderking/learn-pub-sub-starter/internal/gamelogic"
	"github.com/tenderking/learn-pub-sub-starter/internal/pubsub"
	"github.com/tenderking/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		// Correctly acknowledge the message
		return pubsub.Ack

	}
}
func handlerMoves(gs *gamelogic.GameState, ch *amqp091.Channel) func(am gamelogic.ArmyMove) pubsub.Acktype {
	return func(am gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(am)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+am.Player.Username,
				gamelogic.RecognitionOfWar{
					Defender: gs.Player,
					Attacker: am.Player,
				},
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		default:
			return pubsub.NackDiscard
		}
	}
}

func handlerWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(wr gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(wr)
		message := fmt.Sprintf("%s won a war against %s", winner, loser)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			err := pubsub.PublishGob(
				ch,
				routing.ExchangePerilTopic,
				gs.GetUsername(),
				routing.GameLog{
					CurrentTime: time.Now(),
					Message:     message,
					Username:    gs.GetUsername(),
				},
			)
			if err != nil {
				return pubsub.NackRequeue // Or handle the error differently
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			err := pubsub.PublishGob(
				ch,
				routing.ExchangePerilTopic,
				gs.GetUsername(),
				routing.GameLog{
					CurrentTime: time.Now(),
					Message:     message,
					Username:    gs.GetUsername(),
				},
			)
			if err != nil {
				return pubsub.NackRequeue // Or handle the error differently
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			err := pubsub.PublishGob(
				ch,
				routing.ExchangePerilTopic,
				gs.GetUsername(),
				routing.GameLog{
					CurrentTime: time.Now(),
					Message:     "War ended in a draw",
					Username:    gs.GetUsername(),
				},
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Println("Error! Unknown war outcome")
			return pubsub.NackDiscard

		}
	}
}
