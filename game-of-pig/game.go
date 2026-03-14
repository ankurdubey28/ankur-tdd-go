package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(roll())
}

func roll() int {
	return rand.Intn(6) + 1
}
