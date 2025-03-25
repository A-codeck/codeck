//Precisa rodar go get github.com/lib/pq pra pegar o pacote do github
package main

import (
    "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    _ "github.com/lib/pq"
)

type Sprites struct {
	FrontDefault string `json:"front_default"`
}

type Pokemon struct {
	Sprites Sprites `json:"sprites"`
}

const (
    host     = "localhost"
    port     = 5432
    user     = "octavui"
    password = "1234"
    dbname   = "pokemons"
)

func connectToDb() *sql.DB{
    var psqlInfo string = fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
      panic(err)
    }
    err = db.Ping()
    if err != nil {
      panic(err)
    }

    return db
}

func pokemonUrl(name string)string{
    db := connectToDb()
    defer db.Close()
    var query string = "SELECT sprite FROM sprites WHERE name=$1"
    row := db.QueryRow(query, name)

    var pokemonUrl string
    rowerr := row.Scan(&pokemonUrl)
    if rowerr != nil {
        log.Printf("Error scanning row")
        log.Fatal(rowerr)
    }
    return pokemonUrl
}

func GetPokemon(w http.ResponseWriter, req *http.Request) {
    var pokename string = req.URL.Query()["pokename"][0];
    var pokemonUrl string = pokemonUrl(pokename) 
	var pokemon Pokemon
    pokemon.Sprites.FrontDefault = pokemonUrl
	response := map[string]string{
		"pokemonUrl": pokemon.Sprites.FrontDefault,
	}
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error creating JSON response: %v", err)
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
		return
	}
	
	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/", GetPokemon)
	fmt.Println("Server running on :8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal(err)
	}
}
