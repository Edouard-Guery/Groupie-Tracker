package main

import (
	"net/http"
	"strconv"
	"strings"
)

func filtrer(r *http.Request) []Artiste {
	creationMin, _ := strconv.Atoi(r.FormValue("creationMin"))
	creationMax, _ := strconv.Atoi(r.FormValue("creationMax"))
	albumMin, _ := strconv.Atoi(r.FormValue("albumMin"))
	albumMax, _ := strconv.Atoi(r.FormValue("albumMax"))
	membres := r.FormValue("membres")
	lieu := r.FormValue("lieu")

	var resultat []Artiste
	for _, artiste := range artistes {
		album, _ := strconv.Atoi(artiste.FirstAlbum[6:])
		nombre := strconv.Itoa(len(artiste.Members))
		garder := true

		if artiste.CreationDate < creationMin {
			garder = false
		}
		if creationMax != 0 && artiste.CreationDate > creationMax {
			garder = false
		}
		if album < albumMin {
			garder = false
		}
		if albumMax != 0 && album > albumMax {
			garder = false
		}
		if membres != "" && !strings.Contains(membres, nombre) {
			garder = false
		}
		if lieu != "" && !joueA(artiste, lieu) {
			garder = false
		}

		if garder {
			resultat = append(resultat, artiste)
		}
	}
	return resultat
}

func joueA(artiste Artiste, recherche string) bool {
	for lieu := range artiste.Concerts {
		if strings.Contains(lieu, recherche) {
			return true
		}
	}
	return false
}
