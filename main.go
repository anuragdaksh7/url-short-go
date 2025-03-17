package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz")
	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type URLPayload struct {
	URL string `json:"url"`
}

var shortURLs map[string]string = make(map[string]string)

var reverseURLs map[string]string = make(map[string]string)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n");
}

func getShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println(shortURLs)
	fmt.Println(reverseURLs)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var payload URLPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", 400)
		return
	}

	url := payload.URL
	fmt.Println("Recieved URL: ", url)

	if _, exists := shortURLs[url]; exists {
		fmt.Println("URL exists in the map")
		fmt.Fprintf(w, "%s", shortURLs[url])
		return
	} else {
		fmt.Println("URL does not exists in the map")
		var shortened string = generateRandomString(5)

		var _, randomExists = reverseURLs[shortened]
		if randomExists {
			fmt.Println("Regeneating short url")
		}
		
		for {
			if !randomExists {
				break
			}
			shortened = generateRandomString(5)
			_, randomExists = reverseURLs[shortened]
		}
		
		shortURLs[url] = shortened
		reverseURLs[shortened] = url
	}

	fmt.Fprintf(w, "%s", shortURLs[url])
}

func handleTo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	parts := strings.Split(path, "/")

	var id string
	if len(parts) >= 3 && parts[1] == "to" {
		id = parts[2]
		fmt.Printf("Received request for ID: %s\n", id)
		// fmt.Fprintf(w, "Received request for ID: %s", id)
	}

	fmt.Printf("Received request for ID: %s", id)
	var _, url = reverseURLs[id]
	if (url) {
		redirectURL := reverseURLs[id]
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return;
	} else {
		http.Error(w, "404", 404)
	}
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/shortURL", getShortURL)
	http.HandleFunc("/to/", handleTo)

	http.ListenAndServe(":8090", nil)
}

