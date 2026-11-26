package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Faux artistes pour les tests
var exemples = []Artiste{
	{ID: 1, Name: "Queen", Members: []string{"Freddie", "Brian"}, CreationDate: 1970, FirstAlbum: "14-12-1973", Lieux: []string{"osaka-japan", "lausanne-switzerland"}},
	{ID: 5, Name: "Solo", Members: []string{"Bob"}, CreationDate: 2010, FirstAlbum: "01-01-2012", Lieux: []string{"new_york-usa"}},
}

// Pour chaque adresse, le nombre d'artistes attendu
func TestFiltres(t *testing.T) {
	tests := map[string]int{
		"/":                             2,
		"/?creationMin=2000":            1,
		"/?creationMax=1960":            0,
		"/?albumMin=1970&albumMax=1980": 1,
		"/?membres=1":                   1,
		"/?membres=1&membres=2":         2,
		"/?lieu=Osaka,+Japan":           1,
		"/?lieu=new+york+usa":           1,
		"/?lieu=usa":                    1,
		"/?creationMin=abc":             2,
		"/?recherche=QUEEN":             1,
		"/?recherche=freddie":           1,
		"/?recherche=New+York,+USA":     1,
		"/?recherche=14-12-1973":        1,
		"/?recherche=2010":              1,
		"/?recherche=inconnu":           0,
	}
	for adresse, attendu := range tests {
		obtenu := len(filtrer(exemples, httptest.NewRequest("GET", adresse, nil)))
		if obtenu != attendu {
			t.Errorf("%s : %d au lieu de %d", adresse, obtenu, attendu)
		}
	}
}

// Vérifie le code de réponse de chaque page
func TestPages(t *testing.T) {
	stock = exemples
	verifierCode(t, accueil, "/", 200)
	verifierCode(t, accueil, "/inconnu", 404)
	verifierCode(t, pageArtiste, "/artiste?id=1", 200)
	verifierCode(t, pageArtiste, "/artiste?id=99", 404)
	verifierCode(t, pageArtiste, "/artiste?id=abc", 404)

	stock = nil
	verifierCode(t, accueil, "/", 503)
}

func verifierCode(t *testing.T, page http.HandlerFunc, adresse string, attendu int) {
	w := httptest.NewRecorder()
	page(w, httptest.NewRequest("GET", adresse, nil))
	if w.Code != attendu {
		t.Errorf("%s : %d au lieu de %d", adresse, w.Code, attendu)
	}
}

// L'API répond dans le désordre : chaque artiste doit garder ses données
func TestRelierParID(t *testing.T) {
	serveur := fauxServeur(map[string]string{
		"/artists":   `[{"id":1},{"id":2}]`,
		"/locations": `{"index":[{"id":2,"locations":["lyon-france"]},{"id":1,"locations":["paris-france"]}]}`,
		"/dates":     `{"index":[{"id":2,"dates":["02-02-2020"]},{"id":1,"dates":["01-01-2020"]}]}`,
		"/relation":  `{"index":[{"id":2,"datesLocations":{"lyon-france":[]}},{"id":1,"datesLocations":{"paris-france":[]}}]}`,
	})
	defer serveur.Close()
	adresseAPI = serveur.URL + "/"

	artistes, err := chargerAPI()
	if err != nil {
		t.Fatal(err)
	}
	_, existe := artistes[0].Concerts["paris-france"]
	if artistes[0].Lieux[0] != "paris-france" || artistes[0].Dates[0] != "01-01-2020" || !existe {
		t.Error("les données sont mal reliées")
	}
}

// Transforme un lieu en coordonnées
func TestPosition(t *testing.T) {
	serveur := fauxServeur(map[string]string{"/": `[{"lat":"48.85","lon":"2.35"}]`, "/vide": `[]`})
	defer serveur.Close()

	adresseCarte = serveur.URL + "/?q="
	if chercherPosition("paris").Lat != "48.85" {
		t.Error("mauvaise position")
	}
	adresseCarte = serveur.URL + "/vide?q="
	if chercherPosition("inconnu").Lat != "" {
		t.Error("la position devait être vide")
	}
}

func fauxServeur(reponses map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(reponses[r.URL.Path]))
	}))
}

// Simule une API en panne : il doit y avoir une erreur
func TestAPIEnPanne(t *testing.T) {
	serveur := httptest.NewServer(http.NotFoundHandler())
	defer serveur.Close()
	adresseAPI = serveur.URL + "/"

	_, err := chargerAPI()
	if err == nil {
		t.Error("une erreur était attendue")
	}
}
