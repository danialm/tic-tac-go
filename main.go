package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/tncardoso/gocurses"
)

const ENTER = 10
const ESC = 27

var symbols = [4]string{" ", "X", "O", string(176)}

func hasWon(state [3][3]int) interface{} {
	if state[0][0] == state[0][1] && state[0][1] == state[0][2] && state[0][0] != 0 {
		return [3][3]int{state[0], {0, 0, 0}, {0, 0, 0}}
	}
	if state[1][0] == state[1][1] && state[1][1] == state[1][2] && state[1][0] != 0 {
		return [3][3]int{{0, 0, 0}, state[1], {0, 0, 0}}
	}
	if state[2][0] == state[2][1] && state[2][1] == state[2][2] && state[2][0] != 0 {
		return [3][3]int{{0, 0, 0}, {0, 0, 0}, state[2]}
	}
	if state[0][0] == state[1][0] && state[1][0] == state[2][0] && state[0][0] != 0 {
		return [3][3]int{{state[0][0], 0, 0}, {state[1][0], 0, 0}, {state[2][0], 0, 0}}
	}
	if state[0][1] == state[1][1] && state[1][1] == state[2][1] && state[0][1] != 0 {
		return [3][3]int{{0, state[0][1], 0}, {0, state[1][1], 0}, {0, state[2][1], 0}}
	}
	if state[0][2] == state[1][2] && state[1][2] == state[2][2] && state[0][2] != 0 {
		return [3][3]int{{0, 0, state[0][2]}, {0, 0, state[1][2]}, {0, 0, state[2][2]}}
	}
	if state[0][0] == state[1][1] && state[1][1] == state[2][2] && state[0][0] != 0 {
		return [3][3]int{{state[0][0], 0, 0}, {0, state[1][1], 0}, {0, 0, state[2][2]}}
	}
	if state[0][2] == state[1][1] && state[1][1] == state[2][0] && state[0][2] != 0 {
		return [3][3]int{{0, 0, state[0][2]}, {0, state[1][1], 0}, {state[2][0], 0, 0}}
	}
	return false
}

func isTie(state [3][3]int) bool {
	count := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if state[i][j] == 1 || state[i][j] == 2 {
				count++
			}
		}
	}
	return count == 9
}

func printState(state [3]int) string {
	s := []string{" ", symbols[state[0]], " | ", symbols[state[1]], " | ", symbols[state[2]], " "}
	return strings.Join(s, "")
}

func printCanvas(state [3][3]int, players []string, currentPlayerIndex int, scores []int) string {
	currentPlayer := players[currentPlayerIndex]
	var scoreBoard = strings.Join([]string{players[0], ": ", strconv.Itoa(scores[0]), "             ", players[1], ": ", strconv.Itoa(scores[1])}, "")
	var prompt = strings.Join([]string{"It's ", currentPlayer, "'s turn..."}, "")
	s := []string{scoreBoard, prompt, "", printState(state[0]), "___|___|___", printState(state[1]), "___|___|___", printState(state[2]), "   |   |   "}

	return strings.Join(s, "\n")
}

func main() {
	println("Hi, what is your name?")
	var players []string = make([]string, 2)
	var scores []int = make([]int, 2)
	fmt.Scan(&players[0])
	fmt.Printf("Welcome %v!\n", players[0])
	println("Let the other player enter their name...")
	fmt.Scan(&players[1])
	fmt.Printf("Hi %v, welcome!\n", players[1])
	var currentPlayerIndex = 0
	println("Your turn %v", players[currentPlayerIndex])
	var state = [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	var gameIsOver = false

	gocurses.Initscr()
	defer gocurses.End()
	gocurses.Cbreak()
	gocurses.Noecho()
	gocurses.Stdscr.Keypad(true)

	counter := 0
	x := 0
	y := 0
	state[y][x] = 3
	gocurses.Addstr(printCanvas(state, players, currentPlayerIndex, scores))
	gocurses.Refresh()
	for !gameIsOver || counter > 100 {
		counter = counter + 1
		var userInput = gocurses.Stdscr.Getch()

		if userInput == ESC {
			gameIsOver = true
		}
		if userInput == gocurses.KEY_LEFT {
			x--
		}
		if userInput == gocurses.KEY_RIGHT {
			x++
		}
		if userInput == gocurses.KEY_UP {
			y--
		}
		if userInput == gocurses.KEY_DOWN {
			y++
		}
		var actualX = int(math.Mod(math.Abs(float64(x)), 3))
		var actualY = int(math.Mod(math.Abs(float64(y)), 3))
		if userInput == ENTER {
			if state[actualY][actualX] == 3 {
				state[actualY][actualX] = currentPlayerIndex + 1
				if currentPlayerIndex == 0 {
					currentPlayerIndex = 1
				} else {
					currentPlayerIndex = 0
				}
			}
			w := hasWon(state)
			if w != false {
				state = w.([3][3]int)
				scores[currentPlayerIndex]++
				gocurses.Clear()
git				gocurses.Addstr(printCanvas(state, players, currentPlayerIndex, scores))
				gocurses.Refresh()
				time.Sleep(3 * time.Second)
				state = [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
			}

			if isTie(state) {
				state = [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
			}
		} else {
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if state[i][j] == 3 {
						state[i][j] = 0
					}
				}
			}
			if state[actualY][actualX] == 0 {
				state[actualY][actualX] = 3
			}
		}
		gocurses.Clear()
		gocurses.Addstr(printCanvas(state, players, currentPlayerIndex, scores))
		gocurses.Refresh()
	}
}
