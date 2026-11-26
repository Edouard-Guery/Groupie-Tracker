# 🎸 Groupie Tracker - Site web d'artistes et de concerts

Site web développé en Go (Golang) qui affiche des artistes, leurs membres, leurs concerts et leur musique à partir d'une API.

---

## 📖 Présentation
Réalisé en groupe dans le cadre de notre formation à **YNOV Campus**, ce projet valide la création d'un serveur web (`net/http`), la lecture d'une API en JSON (`encoding/json`) et l'utilisation de templates HTML (`html/template`) en Go.

Le but est de récupérer les données de l'API [Groupie Trackers](https://groupietrackers.01edu.ai/api) et de les présenter sur un site clair et simple à utiliser.

**Côté Go :** uniquement la bibliothèque standard (aucun package externe).

### 👥 Équipe
* **Edouard GUERY**
* **Gabriel ANTUNES**

---

## 🚀 Fonctionnalités
* **🎤 Liste des artistes :** Tous les artistes de l'API, avec image, année de création et premier album.
* **🔎 Filtres :** Année de création, premier album, nombre de membres et lieu de concert. Les filtres peuvent être combinés et les choix restent affichés.
* **⌨️ Recherche :** Par artiste, membre, lieu, premier album ou date de création, sans tenir compte des majuscules.
* **📄 Page par artiste :** Membres, concerts, musique et carte sur une seule page.
* **🗺️ Géolocalisation :** Les lieux de concert sont convertis en coordonnées par Go, puis placés sur une carte.
* **🎵 Musique :** Un lecteur Deezer affiche les titres les plus écoutés de l'artiste.
* **⏭️ Navigation :** Artiste précédent et suivant, dans l'ordre de l'API.
* **⚠️ Gestion des erreurs :** Une page dédiée s'affiche si l'adresse ou l'artiste n'existe pas (404) ou si l'API est indisponible (503).
* **🔄 Automatique :** Les données sont gardées en mémoire et actualisées chaque minute. Un nouvel artiste apparaît sans modifier le code.
* **⚡ Raccourcis clavier :** `R` recherche, `H` accueil, `P` artiste précédent, `S` artiste suivant, avec `Alt` (ou `Alt + Maj` sur Firefox).
* **📱 Responsive :** Pensé pour PC, adapté aux tablettes et téléphones.

---

## 📁 Arborescence du Projet

```text
Groupie-Tracker/
├── main.go               # Point d'entrée, serveur et pages
├── api.go                # Lecture et stockage des 4 parties de l'API
├── filtres.go            # Filtres et recherche de la page d'accueil
├── locate.go             # Conversion des lieux en coordonnées
├── main_test.go          # Tests unitaires
├── go.mod                # Initialisation du module Go
├── README.md             # Documentation du projet
└── templates/
    ├── index.html        # Page d'accueil
    ├── artiste.html      # Page d'un artiste
    ├── erreur.html       # Page d'erreur
    └── styles.css        # Feuille de style
```

---

## 🛠️ Architecture Technique
Pour garder un code simple et lisible, la logique Go est découpée en 4 fichiers :
1. **Le serveur (`main.go`) :** Associe chaque adresse à une fonction et affiche les pages.
2. **Les données (`api.go`) :** Lit les 4 parties de l'API (`artists`, `locations`, `dates`, `relation`), les relie grâce à l'`ID` et les garde en mémoire. Si l'API tombe en panne, les dernières données restent disponibles.
3. **Les filtres (`filtres.go`) :** Parcourt les artistes et garde ceux qui correspondent aux filtres et à la recherche.
4. **La géolocalisation (`locate.go`) :** Transforme un lieu en coordonnées avec Nominatim et garde le résultat en mémoire.

### Les adresses du site

| Adresse | Rôle |
|---|---|
| `/` | Accueil, avec la liste, la recherche et les filtres |
| `/artiste?id=1` | Page d'un artiste |
| `/coords?lieu=paris-france` | Coordonnées d'un lieu (JSON) |
| `/styles.css` | Feuille de style |

### Ressources externes
- **API Groupie Trackers** : Les artistes, les lieux et les dates de concerts.
- **Leaflet et OpenStreetMap** : L'affichage de la carte.
- **Nominatim** : La recherche des coordonnées d'un lieu.
- **Deezer** : La recherche de l'artiste et le lecteur de musique.
- **Google Fonts** : Les polices Bricolage Grotesque et Figtree.

---

## 🔍 Focus sur le code

**1. La lecture de l'API**

Les 4 parties sont lues l'une après l'autre. Si l'une échoue, l'erreur est renvoyée et les anciennes données sont gardées :
```go
err1 := lireAPI("artists", &artistes)
err2 := lireAPI("locations", &locations)
err3 := lireAPI("dates", &dates)
err4 := lireAPI("relation", &relations)
```

Une goroutine actualise les données chaque minute, sans bloquer le site :
```go
go actualiser()
```

**2. Les filtres**

Une fonction vérifie si une année est comprise entre deux valeurs. Un maximum à 0 veut dire que le filtre n'est pas utilisé :
```go
return valeur >= min && (max == 0 || valeur <= max)
```

**3. La bannière floutée**

L'image de l'artiste est affichée deux fois : nette devant, et floutée en fond :
```css
.fond {
    filter: blur(40px) brightness(0.4) saturate(1.5);
}
```

---

## 🕹️ Lancement rapide

### Prérequis
* **Go** : Version 1.22 ou supérieure.
* Une connexion internet.
* Un navigateur web moderne (Chrome, Firefox, Safari...).

### Exécuter le programme
À la racine du projet :

```bash
go run .
```
*Le site est ensuite disponible sur http://localhost:8080. Pour arrêter le serveur : `Ctrl+C`.*

### Lancer les tests
Filtres, recherche, pages, erreurs, API et géolocalisation :

```bash
go test -v
```

---

## ⚠️ Limites connues
- L'API ne donne pas l'État des villes : chercher « Washington, USA » ne peut pas trouver un concert à « Seattle, USA ».
- Les repères de la carte apparaissent un par un : Nominatim autorise une seule demande par seconde. Le serveur les garde ensuite en mémoire jusqu'à son arrêt.
- Le lecteur utilise le premier artiste trouvé par Deezer : deux artistes du même nom peuvent être confondus.
- Sans compte Deezer connecté, seuls des extraits peuvent être écoutés.
