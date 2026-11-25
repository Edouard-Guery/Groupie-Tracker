package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func main() {
	chargerArtistes()
	chargerConcerts()

	http.HandleFunc("/", accueil)
	http.HandleFunc("/artiste", pageArtiste)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func accueil(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", 302)
		return
	}

	page, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintln(w, "Erreur : template introuvable")
		return
	}
	page.Execute(w, filtrer(r))
}

func pageArtiste(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	if id < 1 || id > len(artistes) {
		http.Redirect(w, r, "/", 302)
		return
	}

	page, err := template.ParseFiles("templates/artiste.html")
	if err != nil {
		fmt.Fprintln(w, "Erreur : template introuvable")
		return
	}
	page.Execute(w, artistes[id-1])
}
