package main

import (
	"fmt"
)

func main() {
	result := make(chan string)
	ensure := make(chan string)

	go Client(result)
	go Server(ensure)

	r := <-ensure
	fmt.Println(r)
	r = <-result
	fmt.Println(r)
}
