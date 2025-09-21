package main

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

func logCommand(prefix, message string) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, prefix, message)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	flush := func(s string) {
		fmt.Fprintln(writer, s)
		writer.Flush()
	}

	engineColor := NotInitialized
	moveHistory := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		logCommand("IN", line)
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
			moveHistory = []string{}
			engineColor = NotInitialized
			continue

		case "position":
			if engineColor == NotInitialized {
				moveHistory = getUCIMoves(fields)
				if len(moveHistory)%2 == 0 {
					engineColor = White
				} else {
					engineColor = Black
				}
			}
			continue

		case "go":
			if engineColor == White {
				flush("bestmove b2b4")
			} else {
				flush("bestmove b7b5")
			}
		}
	}
}
