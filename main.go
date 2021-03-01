package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates go.mod
var templatesFS embed.FS

func main() {
	allFiles, err := templatesFS.ReadDir(".")
	if err != nil {
		log.Fatalf("Couldn't read dir: %s", err)
	}
	for _, file := range allFiles {
		log.Println(file.Name())
	}
	template, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		log.Fatalf("error parsing templates filesystem: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, struct{}{})
	})

	log.Printf("Running on port 8080")
	http.ListenAndServe(":8080", mux)
}
