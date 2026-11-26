package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Un artiste avec toutes ses informations
type Artiste struct {
	ID           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Lieux        []string
	Dates        []string
	Concerts     map[string][]string
}

// Structures qui reprennent la forme du JSON de l'API
type Locations struct {
	Index []struct {
		ID        int
		Locations []string
	}
}

type Dates struct {
	Index []struct {
		ID    int
		Dates []string
	}
}

type Relations struct {
	Index []struct {
		ID             int
		DatesLocations map[string][]string
	}
}

var adresseAPI = "https://groupietrackers.01edu.ai/api/"

// Artistes gardés en mémoire, et verrou pour ne pas lire et écrire en même temps
var stock []Artiste
var verrou sync.Mutex

// Lit une partie de l'API et range le JSON dans resultat
func lireAPI(partie string, resultat any) error {
	client := http.Client{Timeout: 10 * time.Second}
	reponse, err := client.Get(adresseAPI + partie)
	if err != nil {
		return err
	}
	defer reponse.Body.Close()
	return json.NewDecoder(reponse.Body).Decode(resultat)
}

// Récupère les 4 parties de l'API et les rassemble
func chargerAPI() ([]Artiste, error) {
	var artistes []Artiste
	var locations Locations
	var dates Dates
	var relations Relations

	err1 := lireAPI("artists", &artistes)
	err2 := lireAPI("locations", &locations)
	err3 := lireAPI("dates", &dates)
	err4 := lireAPI("relation", &relations)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return nil, errors.New("l'API ne répond pas")
	}

	// Relie à chaque artiste ses lieux, ses dates et ses concerts grâce à l'ID
	for i := range artistes {
		for _, location := range locations.Index {
			if location.ID == artistes[i].ID {
				artistes[i].Lieux = location.Locations
			}
		}
		for _, date := range dates.Index {
			if date.ID == artistes[i].ID {
				artistes[i].Dates = date.Dates
			}
		}
		for _, relation := range relations.Index {
			if relation.ID == artistes[i].ID {
				artistes[i].Concerts = relation.DatesLocations
			}
		}
	}
	return artistes, nil
}

// Recharge les données toutes les minutes
func actualiser() {
	for {
		artistes, err := chargerAPI()
		if err != nil {
			fmt.Println("Erreur :", err)
		} else {
			verrou.Lock()
			stock = artistes
			verrou.Unlock()
		}
		time.Sleep(time.Minute)
	}
}

// Donne les artistes gardés en mémoire
func lireStock() []Artiste {
	verrou.Lock()
	defer verrou.Unlock()
	return stock
}
