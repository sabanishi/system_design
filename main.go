package main

import "fmt"

func main() {
	greet := hello("Masu")
	fmt.Println(greet)
}

func hello(name string) string{
	return fmt.Sprintf("Hello, %s!",name)
}
