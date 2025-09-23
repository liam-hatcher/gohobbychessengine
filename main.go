package main

import (
	"github.com/liam-hatcher/gohobbyengine/chess"
	"github.com/liam-hatcher/gohobbyengine/uci"
)

func main() {
	position := chess.NewPosition()
	engine := uci.NewEngine()
	engine.Run(position)
}
