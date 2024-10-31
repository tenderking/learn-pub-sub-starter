package main

import (
	"fmt"
	// Import syscall for SIGTERM

	"github.com/tenderking/learn-pub-sub-starter/internal/gamelogic"
	"github.com/tenderking/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}
func handlerMoves(gs *gamelogic.GameState) func(am gamelogic.ArmyMove) {
	return func(am gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(am)
	}
}
