package main

import (
	"encoding/json"
	"net/http"
)

type Artiste struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// recupererArtistes va chercher les artistes
func recupererArtistes() []Artiste {
	reponse, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil
	}
	defer reponse.Body.Close()

	//JSON en liste d'artistes
	var artistes []Artiste
	json.NewDecoder(reponse.Body).Decode(&artistes)
	return artistes
}

// func
