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
	for {
		roll := rollDice()
		if roll == 1 {
			fmt.Println("got 1, your chance over")
			// toggle player
			playerId = 1 - playerId
		}
		fmt.Printf("PlayerId:%d Would you like to hold?", playerId)
		var hold string
		fmt.Scan(&hold)
		if strings.ToLower(hold) == "yes" {
			// add rolled value to total score of player
			fmt.Println("player score updated")
			SCORE[playerId] += roll

			// toggle player
			playerId = 1 - playerId
		}
		// check if any winner
		if playerId, ok := findWinner(SCORE); ok {
			fmt.Printf("winner found , player ID: %d", playerId)
		}
		fmt.Println(SCORE)
		SCORE[playerId] += roll
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
