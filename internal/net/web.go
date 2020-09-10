package net

import (
	"html/template"
	"log"
	"net/http"
)

// GET request Webpage serving
func web(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("mediaserver.html")
		err := t.Execute(w, nil)
		if err != nil {
			log.Println("[ERROR] Template Execute : ", err)
		}
	}
}
