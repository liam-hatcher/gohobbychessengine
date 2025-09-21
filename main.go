package main

import (
	"github.com/liam-hatcher/gohobbyengine/uci"
)

func main() {
	engine := uci.NewEngine()
	engine.Run()
}
