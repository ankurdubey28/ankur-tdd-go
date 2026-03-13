package main

func main() {

}

func IsEven(arr []int) []int {
	nums := make([]int, 0)
	for _, v := range arr {
		if v%2 == 0 {
			nums = append(nums, v)
		}
	}
	return nums
}

func IsOdd(arr []int) []int {
	nums := make([]int, 0)
	for _, v := range arr {
		if v%2 != 0 {
			nums = append(nums, v)
		}
	}
	return nums
}

func IsPrime(arr []int) []int {
	var nums []int

	for _, v := range arr {
		if v == 0 || v == 1 {
			continue
		}
		found := false
		for i := 2; i*i <= v; i++ {
			if v%i == 0 {
				found = true
				break
			}
		}
		if !found {
			nums = append(nums, v)
		}
	}
	return nums
}

func IsPrimeOdd(arr []int) []int {
	primes := IsPrime(arr)
	OddPrimes := IsOdd(primes)
	return OddPrimes
}

func IsEven5(arr []int) []int {
	even := IsEven(arr)
	var nums []int
	for _, v := range even {
		if v%5 == 0 {
			nums = append(nums, v)
		}
	}
	return nums
}

func IsOdd3(arr []int) []int {
	odd := IsOdd(arr)
	var nums []int

	for _, v := range odd {
		if v%3 == 0 && v > 10 {
			nums = append(nums, v)
		}
	}
	return nums
}

type Predicate[T any] func(T) bool

// Predicate builders

// Generator
func GreaterThan(n int) Predicate[int] {
	return func(i int) bool {
		return n < i
	}
}

func LesserThan(n int) Predicate[int] {
	return func(i int) bool {
		return n > i
	}
}

func Even() func(int) bool {
	return func(i int) bool {
		return i%2 == 0
	}
}

func Odd() func(int) bool {
	return func(i int) bool {
		return i%2 != 0
	}
}

func Prime() func(int) bool {
	return func(x int) bool {
		if x == 0 || x == 1 {
			return false
		}
		for i := 2; i*i <= x; i++ {
			if x%i == 0 {
				return false
			}
		}
		return true
	}
}

func Multiple(n int) Predicate[int] {
	return func(i int) bool {
		return i%n == 0
	}
}

// all
func All[T any](preds ...Predicate[T]) Predicate[T] {
	return func(v T) bool {
		for _, p := range preds {
			if !p(v) {
				return false
			}
		}
		return true
	}
}

// any
func Any[T any](preds ...Predicate[T]) Predicate[T] {
	return func(v T) bool {
		for _, p := range preds {
			if p(v) {
				return true
			}
		}
		return false
	}
}

func Filter[T any](items []T, pred Predicate[T]) []T {
	var res []T
	for _, v := range items {
		if pred(v) {
			res = append(res, v)
		}
	}
	return res
}
