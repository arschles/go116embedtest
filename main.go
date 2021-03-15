package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates go.mod README.md
var templatesFS embed.FS

//go:embed README.md
var readmeFS embed.FS

type catResponse struct {
	URL string `json:"url"`
}

type dogResponse struct {
	URL string `json:"message"`
}

func main() {
	template, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		log.Fatalf("error parsing templates filesystem: %s", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/cat", func(w http.ResponseWriter, r *http.Request) {
		catResp, err := http.Get("https://api.thecatapi.com/v1/images/search")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error getting cat! %s", err)))
			return
		}
		defer catResp.Body.Close()
		respList := []catResponse{}
		if err := json.NewDecoder(catResp.Body).Decode(&respList); err != nil {
			w.Write([]byte(fmt.Sprintf("Error in cat response JSON: %s", err)))
			return
		}
		catURL := respList[0].URL
		template.Execute(w, map[string]string{
			"AnimalURL": catURL,
		})
	})

	mux.HandleFunc("/dog", func(w http.ResponseWriter, r *http.Request) {
		dogResp, err := http.Get("https://dog.ceo/api/breeds/image/random")
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error getting dog! %s", err)))
			return
		}
		defer dogResp.Body.Close()
		respList := dogResponse{}
		if err := json.NewDecoder(dogResp.Body).Decode(&respList); err != nil {
			w.Write([]byte(fmt.Sprintf("Error in dog response JSON: %s", err)))
			return
		}
		dogURL := respList.URL
		template.Execute(w, map[string]string{
			"AnimalURL": dogURL,
		})
	})

	mux.HandleFunc("/readme", func(w http.ResponseWriter, r *http.Request) {
		readmeBytes, err := readmeFS.ReadFile("README.md")
		if err != nil {
			log.Fatalf("error reading README: %s", err)
		}
		w.Write(readmeBytes)
	})

	log.Printf("Running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
