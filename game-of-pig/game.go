package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var SCORE = map[int]int{
	0: 0,
	1: 0,
}

func main() {
	playerId := getRandomId()
	turnTotal := 0
	for {
		roll := rollDice()
		fmt.Printf("Rolled and got: %v", roll)
		turnTotal += roll
		fmt.Println()
		if roll == 1 {
			fmt.Println("got 1, your chance over")
			// toggle player && turnTotal=0
			turnTotal = 0
			playerId = 1 - playerId
			continue
		}
		fmt.Printf("PlayerId:%d Would you like to hold?", playerId)
		var hold string
		fmt.Scan(&hold)
		if strings.ToLower(hold) == "yes" {
			// add turnTotal value to total score of player
			fmt.Println("player score updated")
			SCORE[playerId] += turnTotal
			// toggle player
			playerId = 1 - playerId
			turnTotal = 0
		}
		// check if any winner
		if pID, ok := findWinner(SCORE); ok {
			fmt.Printf("winner found , player ID: %d", pID)
		}
		fmt.Println(turnTotal)
	}
}

func rollDice() int {
	return rand.Intn(6) + 1
}

func getRandomId() int {
	return rand.Intn(2)
}

func findWinner(score map[int]int) (int, bool) {
	for k, v := range score {
		if v >= 10 {
			return k, true
		}
	}
	return -1, false
}
