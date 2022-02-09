package main

import (
	"fmt"
)

func main() {
	fmt.Println("Returned:", MyFunc())
}

func MyFunc() (ret bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic occurred:", fmt.Sprint(r))
		}
	}()
	ch := make(chan int, 10)
	close(ch)
	ch <- 1
	return true
}
