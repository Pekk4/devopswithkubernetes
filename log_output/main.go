package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func getPongs() string {
	resp, err := http.Get("http://pingpong-service:80/pings")
	if err != nil {
		return "Could not get pongs: " + err.Error()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Could not read response: " + err.Error()
	}

	return "Ping / pongs: " + string(body)
}

func main() {
	s := uuid.New().String()

	role := os.Getenv("ROLE")
	if role == "writer" {
		log.Println("Started with writer mode, writing to 'output.log' and not starting a server...")

		f, err := os.OpenFile("/logs/output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			for {
				timestamp := time.Now().Format("2006/01/02 15:04:05")
				f.WriteString("[" + timestamp + "] " + s + "\n")
				time.Sleep(5 * time.Second)
			}
		}
	} else {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}

		log.Println("Server started in port " + port)

		http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			log_data, err := os.ReadFile("/logs/output.log")
			if err != nil {
				http.Error(w, "Could not read a log file", http.StatusInternalServerError)
				return
			}
			pong_data := getPongs()
			w.Write(log_data)
			w.Write([]byte(pong_data))
		})

		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
