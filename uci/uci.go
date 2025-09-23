package uci

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/liam-hatcher/gohobbyengine/chess"
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

func (e *Engine) HandleGo(p *chess.Position) string {
	if e.EngineColor == White {
		// play some random pawn moves for now
		moves := p.GenerateWhitePawnMoves()
		uciMoves := make([]string, len(moves))
		for i, m := range moves {
			uciMoves[i] = chess.ToUCINotation(m)
		}

		LogCommand("DEBUG", fmt.Sprintf("%+v\n", uciMoves))

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomIndex := r.Intn(len(uciMoves))

		randMove := uciMoves[randomIndex]
		p.ApplyMove(uciMoves[randomIndex])

		LogCommand("BLACK Pawns: ", fmt.Sprintf("%d", p.BlackPawns))
		return randMove
	} else {
		// dummy move for black right now=
		return "h7h6"
	}
}

func (e *Engine) Run(p *chess.Position) {
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
			e.MoveHistory = getUCIMoves(fields)
			if e.EngineColor == NotInitialized {
				if len(e.MoveHistory)%2 == 0 {
					e.EngineColor = White
				} else {
					e.EngineColor = Black
				}
			}
			if len(e.MoveHistory) > 0 {
				lastMove := e.MoveHistory[len(e.MoveHistory)-1]
				p.ApplyMove(lastMove)
			}
		case "go":
			move := e.HandleGo(p)
			if move != "" {
				flush("bestmove " + move)
			}
		}
	}
}
