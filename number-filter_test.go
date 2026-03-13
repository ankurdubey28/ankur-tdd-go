package main

import (
	"slices"
	"testing"
)

func TestEven(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{2, 4, 6, 8, 10}
	got := IsEven(input)
	if !slices.Equal(want, got) {
		t.Errorf("IsEven(%v) = %v; want %v", input, got, want)
	}
}

func TestOdd(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{1, 3, 5, 7, 9}
	got := IsOdd(input)
	if !slices.Equal(got, want) {
		t.Errorf("IsEven(%v) = %v; want %v", input, got, want)
	}
}

func TestPrime(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{2, 3, 5, 7}
	got := IsPrime(input)
	if !slices.Equal(want, got) {
		t.Errorf("IsPrime(%v)=%v; want %v", input, got, want)
	}
}

func TestPrimeOdd(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{3, 5, 7}
	got := IsPrimeOdd(input)
	if !slices.Equal(want, got) {
		t.Errorf("IsPrimeOdd(%v)==%v , want %v", input, got, want)
	}
}

func TestEven5(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	want := []int{10, 20}
	got := IsEven5(input)
	if !slices.Equal(got, want) {
		t.Errorf("IsPrime5(%v)==%v, want %v", input, got, want)
	}
}

func TestOdd3(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	want := []int{15}
	got := IsOdd3(input)
	if !slices.Equal(want, got) {
		t.Errorf("IsOdd3(%v)==%v , want %v", input, got, want)
	}
}

func TestAll(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	want := []int{9, 15}
	cond := All[int](
		Odd(),
		GreaterThan(5),
		Multiple(3),
	)
	got := Filter(input, cond)
	if !slices.Equal(want, got) {
		t.Errorf("All(%v)==%v , want %v", input, got, want)
	}
}

func TestAny(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	want := []int{2, 3, 5, 7, 10, 11, 13, 15, 16, 17, 18, 19, 20}
	cond := Any[int](
		Prime(),
		GreaterThan(15),
		Multiple(5),
	)
	got := Filter(input, cond)
	if !slices.Equal(got, want) {
		t.Errorf("Any(%v)=%v , want %v", input, got, want)
	}
}
