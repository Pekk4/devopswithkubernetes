package main

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func main() {
	s := uuid.New().String()

	for {
		log.Println(s)
		time.Sleep(5 * time.Second)
	}
}
