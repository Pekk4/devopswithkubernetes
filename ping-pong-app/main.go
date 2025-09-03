package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	mutex sync.Mutex
	count int
)

func pingPongHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Fprintf(w, "pong %d", count)

	count++
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/pingpong", pingPongHandler)

	log.Println("Server started in port" + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
