package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	mutex    sync.Mutex
	count    int
	port     string
	db_creds string
	ctx      context.Context
	store    *CounterStore
)

func pingPongHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Fprintf(w, "pong %d", count)

	count++

	_, err := store.IncrementCounts(ctx)
	if err != nil {
		log.Printf("failed to increment counts: %v", err)
	}

	f, err := os.OpenFile("/logs/pong.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err == nil {
		defer f.Close()
		fmt.Fprintf(f, "%d\n", count)
	}
}

func initSession() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	db_creds = os.Getenv("PG_URL")
	if db_creds == "" {
		log.Fatal("variable PG_URL is not set")
	}

	var err error
	store, err = NewCounterStore(db_creds)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	ctx = context.Background()
	if err := store.Init(ctx); err != nil {
		log.Fatalf("failed to initialize DB: %v", err)
	}
}

func main() {
	initSession()
	defer func() {
		if store != nil {
			_ = store.Close()
		}
	}()

	http.HandleFunc("/pingpong", pingPongHandler)
	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "%d", count)
		result, err := store.GetCounts(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get counts: %result", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%d", result)
	})

	log.Println("Server started in port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
