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

func getLogs() string {
	data, err := os.ReadFile("/logs/output.log")
	if err != nil {
		return "Could not read a log file: " + err.Error()
	}

	return string(data)
}

func getFileContent() string {
	data, err := os.ReadFile("/information/information.txt")
	if err != nil {
		return "Could not read a log file: " + err.Error()
	}

	return "File content: " + string(data)
}

func main() {
	s := uuid.New().String()

	role := os.Getenv("ROLE")
	if role == "writer" {
		log.Println("Started with writer mode, writing to 'output.log' and not starting a server...")

		for {
			f, err := os.OpenFile("/logs/output.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err == nil {
				timestamp := time.Now().Format("2006/01/02 15:04:05")
				f.WriteString("[" + timestamp + "] " + s + "\n")
				f.Close()
			}
			time.Sleep(5 * time.Second)
		}
	} else {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		msg := os.Getenv("MESSAGE")
		if msg == "" {
			log.Fatal("Environment variable MESSAGE is required")
		}

		log.Println("Server started in port " + port)

		http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			file_data := getFileContent()
			log_data := getLogs()
			pong_data := getPongs()
			w.Write([]byte(file_data))
			w.Write([]byte("env var: " + msg + "\n"))
			w.Write([]byte(log_data))
			w.Write([]byte(pong_data))
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Frostyyy, the snowman..."))
		})

		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
