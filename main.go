package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

/* 
     d9736711 --> {
	                ID: "d9736711",
					OriginalURL: "https://github.com/ankitsingh0805/",
					ShortURL: "d9736711",
					CreationDate: time.Now()}   
*/

var urlDB = make(map[string]URL)

func generateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	data := hasher.Sum(nil)
	fmt.Println("hasher:", data)
	hash := hex.EncodeToString(data)
	fmt.Println("hash:", hash)
	fmt.Println("final string:", hash[:8])
	return hash[:8] // Taking first 8 characters of the hash as short URL
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET method")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	shortURL := createURL(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	fmt.Println("Starting URL shortener...")
	originalURL := "https://github.com/ankitsingh0805/"
	shortURL := createURL(originalURL)
	fmt.Println("Short URL:", shortURL)

	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", ShortURLHandler)

	fmt.Println("Starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}


