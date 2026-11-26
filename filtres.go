package main

import (
	"net/http"
	"strconv"
	"strings"
)

// Garde les artistes qui respectent tous les filtres
func filtrer(artistes []Artiste, r *http.Request) []Artiste {
	// Récupère les choix du visiteur
	recherche := r.FormValue("recherche")
	lieu := nettoyer(r.FormValue("lieu"))
	membres := strings.Join(r.Form["membres"], "")

	var resultat []Artiste
	for _, artiste := range artistes {
		// Prépare l'année du premier album, le nombre de membres et les lieux
		date := strings.Split(artiste.FirstAlbum, "-")
		album, _ := strconv.Atoi(date[len(date)-1])
		nombre := strconv.Itoa(min(len(artiste.Members), 8))
		lieux := nettoyer(strings.Join(artiste.Lieux, " ; "))

		// L'artiste est gardé si toutes les conditions sont vraies
		if entre(artiste.CreationDate, r.FormValue("creationMin"), r.FormValue("creationMax")) &&
			entre(album, r.FormValue("albumMin"), r.FormValue("albumMax")) &&
			(membres == "" || strings.Contains(membres, nombre)) &&
			strings.Contains(lieux, lieu) &&
			correspond(artiste, recherche) {
			resultat = append(resultat, artiste)
		}
	}
	return resultat
}

// Cherche le texte dans toutes les informations de l'artiste
func correspond(artiste Artiste, recherche string) bool {
	texte := artiste.Name + ";" + strconv.Itoa(artiste.CreationDate) + ";" + artiste.FirstAlbum + ";"
	texte += strings.Join(artiste.Members, ";") + ";" + strings.Join(artiste.Lieux, ";")
	return strings.Contains(nettoyer(texte), strings.TrimSpace(nettoyer(recherche)))
}

// Vérifie qu'une valeur est entre un minimum et un maximum (0 = pas de maximum)
func entre(valeur int, texteMin string, texteMax string) bool {
	min, _ := strconv.Atoi(texteMin)
	max, _ := strconv.Atoi(texteMax)
	return valeur >= min && (max == 0 || valeur <= max)
}

func nettoyer(texte string) string {
	texte = strings.ToLower(texte)
	texte = strings.ReplaceAll(texte, "_", " ")
	texte = strings.ReplaceAll(texte, "-", " ")
	texte = strings.ReplaceAll(texte, ",", " ")
	return " " + strings.Join(strings.Fields(texte), " ")
}
