package main

import (
	"fmt"
	"math/rand"
)

type strategy func(int) bool

var WINS = map[int]int{
	0: 0,
	1: 0,
}

var STRATEGY = map[string]strategy{
	"hold least 10": strategyGenerator(10),
	"hold least 25": strategyGenerator(25),
}

func main() {
	// number of games to played
	k := 1000
	// curr strategy map
	var currStrategy = map[int]strategy{
		0: STRATEGY["hold least 10"],
		1: STRATEGY["hold least 25"],
	}
	for k > 0 {
		var score = map[int]int{
			0: 0,
			1: 0,
		}
		playGame(score, currStrategy)
		k -= 1
	}
	fmt.Println(WINS)
}

func playGame(score map[int]int, currStrategy map[int]strategy) {
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
		hold := currStrategy[playerId](turnTotal)
		if hold {
			// add turnTotal value to total score of player
			fmt.Println("player score updated")
			score[playerId] += turnTotal
			// toggle player
			playerId = 1 - playerId
			turnTotal = 0
		}
		// check if any winner
		if pID, ok := findWinner(score); ok {
			fmt.Printf("winner found , player ID: %d", pID)
			// winner found , increment wins for winning player and ,Reset score(happens automatically) and break
			WINS[pID] += 1
			break
		}
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
		if v >= 100 {
			return k, true
		}
	}
	return -1, false
}

func strategyGenerator(k int) strategy {
	return func(turnTotal int) bool {
		return turnTotal >= k
	}
}
