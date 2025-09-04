package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	latestImageTimestamp time.Time
	imageMutex           sync.Mutex
)

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

func main() {
	go handleImageProcedure()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started in port " + port)

	fs := http.FileServer(http.Dir("/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		timestamp := latestImageTimestamp.Format("20060102150405")
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, map[string]string{"ImageTS": timestamp})

		go handleImageProcedure()
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
