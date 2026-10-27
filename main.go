package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", accueil)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Message homepage
func accueil(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur Groupie Tracker !")

	// Affichage noms des artistes
	artistes := recupererArtistes()
	for _, artiste := range artistes {
		fmt.Fprintln(w, artiste.Name)
	}
}

//
