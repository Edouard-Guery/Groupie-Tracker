package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Coordonnées d'un lieu
type Position struct {
	Lat string
	Lon string
}

// Positions déjà trouvées, et verrou pour faire une seule recherche à la fois
var positions = map[string]Position{}
var verrouCarte sync.Mutex

var adresseCarte = "https://nominatim.openstreetmap.org/search?format=json&limit=1&q="

// Envoie la position d'un lieu, par exemple /coords?lieu=paris-france
func coords(w http.ResponseWriter, r *http.Request) {
	verrouCarte.Lock()
	defer verrouCarte.Unlock()

	// "osaka-japan" devient "osaka, japan"
	lieu := strings.ReplaceAll(r.FormValue("lieu"), "_", " ")
	lieu = strings.ReplaceAll(lieu, "-", ", ")

	if positions[lieu].Lat == "" {
		positions[lieu] = chercherPosition(lieu)
	}
	json.NewEncoder(w).Encode(positions[lieu])
}

// Demande la position à Nominatim (OpenStreetMap)
func chercherPosition(lieu string) Position {
	// Nominatim autorise une demande par seconde
	time.Sleep(time.Second)

	// Nominatim demande le nom de l'application
	demande, err := http.NewRequest("GET", adresseCarte+url.QueryEscape(lieu), nil)
	if err != nil {
		return Position{}
	}
	demande.Header.Set("User-Agent", "groupie-tracker-ynov")

	client := http.Client{Timeout: 10 * time.Second}
	reponse, err := client.Do(demande)
	if err != nil {
		fmt.Println("Erreur carte :", err)
		return Position{}
	}
	defer reponse.Body.Close()

	// Garde le premier résultat
	var resultats []Position
	json.NewDecoder(reponse.Body).Decode(&resultats)
	if len(resultats) == 0 {
		return Position{}
	}
	return resultats[0]
}
