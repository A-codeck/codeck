package main
import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "io/ioutil"
)

func GetPokemon(w http.ResponseWriter, req *http.Request){
    var name string = req.URL.Query()["pokename"][0];
    resp,err := http.Get("https://pokeapi.co/api/v2/pokemon/"+name)
    if err != nil {
        log.Fatal(err)
        return
    }
    if resp.Body != nil {
		defer resp.Body.Close()
	}
    body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    

    sprites := result["sprites"].(map[string]interface{})
    fmt.Println(sprites["front_default"]);
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"pokemonUrl":"`+sprites["front_default"].(string)+`"}`))
}
func main(){
    http.HandleFunc("/", GetPokemon)
    http.ListenAndServe(":8001", nil)
}
