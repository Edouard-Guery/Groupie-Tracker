package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

// Données envoyées à la page d'accueil
type PageAccueil struct {
	Artistes []Artiste
	Choix    url.Values
	Coches   map[string]bool
}

// Données envoyées à la page d'un artiste
type PageArtiste struct {
	Artiste   Artiste
	Precedent Artiste
	Suivant   Artiste
}

func main() {
	// Charge les données en arrière-plan
	go actualiser()

	// Chaque adresse appelle une fonction
	http.HandleFunc("/", accueil)
	http.HandleFunc("/artiste", pageArtiste)
	http.HandleFunc("/coords", coords)
	http.HandleFunc("/styles.css", styles)

	// Lance le serveur
	fmt.Println("http://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}

// Page d'accueil : liste, recherche et filtres
func accueil(w http.ResponseWriter, r *http.Request) {
	artistes := lireStock()
	if r.URL.Path != "/" {
		afficher(w, 404, "erreur.html", "Cette page n'existe pas.")
		return
	}
	if len(artistes) == 0 {
		afficher(w, 503, "erreur.html", "Données indisponibles, réessayez dans une minute.")
		return
	}

	// Prépare les artistes et garde les choix du visiteur
	var donnees PageAccueil
	donnees.Artistes = filtrer(artistes, r)
	donnees.Choix = r.Form
	donnees.Coches = map[string]bool{}
	for _, nombre := range r.Form["membres"] {
		donnees.Coches[nombre] = true
	}
	afficher(w, 200, "index.html", donnees)
}

// Page d'un artiste, par exemple /artiste?id=1
func pageArtiste(w http.ResponseWriter, r *http.Request) {
	artistes := lireStock()
	id, _ := strconv.Atoi(r.FormValue("id"))

	for i := range artistes {
		if artistes[i].ID == id {
			var donnees PageArtiste
			donnees.Artiste = artistes[i]
			// Artiste précédent et suivant, s'ils existent
			if i > 0 {
				donnees.Precedent = artistes[i-1]
			}
			if i < len(artistes)-1 {
				donnees.Suivant = artistes[i+1]
			}
			afficher(w, 200, "artiste.html", donnees)
			return
		}
	}
	afficher(w, 404, "erreur.html", "Artiste introuvable.")
}

// Envoie le fichier CSS
func styles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/styles.css")
}

// Remplit un fichier HTML avec les données et l'envoie
func afficher(w http.ResponseWriter, code int, fichier string, donnees any) {
	page, err := template.ParseFiles("templates/" + fichier)
	if err != nil {
		http.Error(w, "Erreur 500 : page indisponible", 500)
		return
	}
	w.WriteHeader(code)
	err = page.Execute(w, donnees)
	if err != nil {
		fmt.Println("Erreur :", err)
	}
}
