package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	latestImageTimestamp time.Time
	imageMutex           sync.Mutex
	port                 string
	backendBaseURL       string
)

type TodoItem struct {
	Todo string `json:"todo"`
	ID   int    `json:"id"`
	Done bool   `json:"done"`
}

func loadConfigFromENV() {
	port = os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}
	backendBaseURL = os.Getenv("BACKEND_URL")
	if backendBaseURL == "" {
		log.Fatal("BACKEND_URL environment variable is required")
	}
}

func getImage() error {
	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	latestImageTimestamp = time.Now().UTC()
	timestamp := latestImageTimestamp.Format("20060102150405")
	imageName := "image_" + timestamp + ".jpg"

	log.Println("Setting timestamp to", timestamp)

	img, err := os.Create("/static/" + imageName)
	if err != nil {
		return err
	}
	defer img.Close()

	logFile, err := os.OpenFile("/logs/image.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err == nil {
		defer logFile.Close()
		logFile.WriteString(timestamp)
	}

	_, err = img.ReadFrom(resp.Body)
	return err
}

func solveLatestImageStatus() {
	data, err := os.ReadFile("/logs/image.log")
	if err == nil {
		timestampStr := string(data)

		t, err := time.ParseInLocation("20060102150405", timestampStr, time.UTC)
		if err == nil {
			latestImageTimestamp = t
			log.Println("Found latest image timestamp:", timestampStr)
		}
	}
}

func handleImageProcedure() {
	imageMutex.Lock()
	defer imageMutex.Unlock()

	if latestImageTimestamp.IsZero() {
		solveLatestImageStatus()
		if latestImageTimestamp.IsZero() {
			err := getImage()
			if err != nil {
				log.Println("Failed to fetch image:", err)
			}
		} else {
			if time.Since(latestImageTimestamp) > 10*time.Minute {
				err := getImage()
				if err != nil {
					log.Println("Failed to refresh image:", err)
				}
			}
		}
	} else {
		if time.Since(latestImageTimestamp) > 10*time.Minute {
			err := getImage()
			if err != nil {
				log.Println("Failed to refresh image:", err)
			}
		}
	}
}

func getTodos() ([]TodoItem, error) {
	resp, err := http.Get(backendBaseURL + "/todos")
	if err != nil {
		log.Println("Failed to fetch todos:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read todos response:", err)
		return nil, err
	}

	var result struct {
		Todos []TodoItem `json:"todos"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Failed to parse todos JSON:", err)
		return nil, err
	}

	return result.Todos, nil
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	resp, err := http.Get(backendBaseURL + "/healthz")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	loadConfigFromENV()

	go handleImageProcedure()

	log.Println("Server started in port " + port)

	fs := http.FileServer(http.Dir("/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		todos, err := getTodos()
		if err != nil {
			log.Println("Error fetching todos:", err)
		}

		timestamp := latestImageTimestamp.Format("20060102150405")
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, map[string]any{
			"ImageTS": timestamp,
			"Todos":   todos,
		})

		go handleImageProcedure()
	})
	http.HandleFunc("/healthz", healthzHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
