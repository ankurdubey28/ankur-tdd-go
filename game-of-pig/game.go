package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(rollDice())
}

func rollDice() int {
	return rand.Intn(6) + 1
}
