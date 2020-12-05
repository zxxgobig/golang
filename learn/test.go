package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("main start")
	go func() {
		log.Fatal("dd")
	}()

	select {}
}
