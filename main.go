package main

import "fmt"

func main() {
	for i := 1; i <= 30; i++ {
		fmt.Println("%d: %s", i, doFizzBuzz(i))
	}
}

func isTree(i int) bool {
	return i%3 == 0
}

func isFive(i int) bool {
	return i%5 == 0
}

func isFifteen(i int) bool {
	return i%15 == 0
}

func doFizzBuzz(i int) string {
	switch {
	case isFifteen(i):
		return "FizzBuzz"
	case isTree(i):
		return "Fizz"
	case isFive(i):
		return "Buzz"
	default:
		return fmt.Sprintf("%d", i)
	}
}
