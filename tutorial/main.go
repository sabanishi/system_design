package main

import "fmt"

func main() {
	var a1 [5]int
	a2 := [...]int{0, 1, 2, 3, 4}
	fmt.Println(a1[0])
	fmt.Println(a2[1])

	//v := 0
	//
	//v = v + 1
	//fmt.Printf("v = %d\n", v)
	//
	//p := &v
	//*p = *p + 1
	//fmt.Printf("v = %d\n", v)

	//for i := 1; i <= 30; i++ {
	//	fmt.Println("%d: %s", i, doFizzBuzz(i))
	//}
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
