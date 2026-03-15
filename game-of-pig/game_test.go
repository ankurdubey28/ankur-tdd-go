package main

import "testing"

func TestRollDice(t *testing.T) {
	got := rollDice()
	want := (got >= 1) && (got <= 6)
	if !want {
		t.Errorf("Wrong dice rolling behaviour ; got %v", got)
	}
}
