package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Artiste struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Concerts     map[string][]string
}

type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relations struct {
	Index []Relation `json:"index"`
}

var artistes []Artiste
var relations Relations

func chargerArtistes() {
	reponse, err := http.Get("https://groupietrackers.01edu.ai/api/artists")
	if err != nil {
		fmt.Println("Erreur API")
		return
	}
	json.NewDecoder(reponse.Body).Decode(&artistes)
	reponse.Body.Close()
}

func chargerConcerts() {
	reponse, err := http.Get("https://groupietrackers.01edu.ai/api/relation")
	if err != nil {
		fmt.Println("Erreur API")
		return
	}
	json.NewDecoder(reponse.Body).Decode(&relations)
	reponse.Body.Close()

	for i, relation := range relations.Index {
		artistes[i].Concerts = relation.DatesLocations
	}
}
