// Gera o commando pra inserir pokemons na DB
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Sprites struct {
	FrontDefault string `json:"front_default"`
}

type Pokemon struct {
	Sprites Sprites `json:"sprites"`
}

func pokemonUrl(name string) string{
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil {
		log.Printf("Error fetching Pokemon data: %v", err)
		return "none"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "none"
	}

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return "none"
	}

    return fmt.Sprintf("('%s', '%s'),", name, pokemon.Sprites.FrontDefault)
}

func main() {
    pokemons := []string{"ditto", "pikachu", "charizard", "blastoise", "voltorb"}
    var command string = "INSERT INTO sprite VALUES "
    for _, pokeName := range pokemons {
        command += pokemonUrl(pokeName) 
    }
    command = command[:len(command)-1]
    command += ";"
    fmt.Println(command)
}

