package uci

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	White = iota
	Black
	NotInitialized
)

type Engine struct {
	MoveHistory   []string
	EngineColor   int
	FirstMoveDone bool
}

func NewEngine() *Engine {
	return &Engine{
		MoveHistory:   []string{},
		EngineColor:   NotInitialized,
		FirstMoveDone: false,
	}
}

// gets the list of moves from a UCI "position" command, e.g. 'position startpos moves e2e4 e7e5'
func getUCIMoves(fields []string) []string {
	var moves []string

	for i, f := range fields {
		if f == "moves" && i+1 < len(fields) {
			moves = fields[i+1:]
			break
		}
	}
	return moves
}

func LogCommand(prefix, message string) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, prefix, message)
}

func (e *Engine) HandleGo() string {
	if e.EngineColor == White {
		return "b2b3"
	} else {
		return "h7h6"
	}
}

func (e *Engine) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	flush := func(s string) {
		fmt.Fprintln(writer, s)
		writer.Flush()
		LogCommand("OUT", s)
	}

	for scanner.Scan() {
		line := scanner.Text()

		LogCommand("IN", line)

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "uci":
			flush("id name GoHobbyEngine")
			flush("id author Liam Hatcher")
			flush("uciok")
		case "isready":
			flush("readyok")
		case "ucinewgame":
			e.MoveHistory = []string{}
			e.FirstMoveDone = false
			e.EngineColor = NotInitialized
		case "position":
			if e.EngineColor == NotInitialized {
				e.MoveHistory = getUCIMoves(fields)
				if len(e.MoveHistory)%2 == 0 {
					e.EngineColor = White
				} else {
					e.EngineColor = Black
				}
			}
		case "go":
			move := e.HandleGo()
			if move != "" {
				flush("bestmove " + move)
			}
		}
	}
}
