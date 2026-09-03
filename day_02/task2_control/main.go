package main

import "fmt"

func main() {

	//FizzBuzz
	i := 1
	for i < 101 {
		switch {
		case i%3 == 0 && i%5 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
		i++
	}
}
