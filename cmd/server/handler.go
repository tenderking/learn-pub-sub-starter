package main

import (
	"fmt"

	"github.com/tenderking/learn-pub-sub-starter/internal/gamelogic"
	"github.com/tenderking/learn-pub-sub-starter/internal/pubsub"
	"github.com/tenderking/learn-pub-sub-starter/internal/routing"
)

func handlerGameLogs() func(routing.GameLog) pubsub.Acktype {
	return func(gl routing.GameLog) pubsub.Acktype {
		fmt.Println("handlerGameLogs called with:", gl)
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(gl)
		if err != nil {
			return pubsub.NackRequeue
		}

		// Correctly acknowledge the message
		return pubsub.Ack

	}
}
