package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", accueil)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func accueil(w http.ResponseWriter, r *http.Request) {
	page, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur 500 : template introuvable", http.StatusInternalServerError)
		return
	}

	// on recup les artistes et on les envoie au template
	artistes := recupererArtistes()
	page.Execute(w, artistes)
}
