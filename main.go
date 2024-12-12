package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

func main() {

	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.HandleFunc("/all/", handlerAllPokemon)
	http.HandleFunc("/pokemon/", handlerPokemon)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handlerPokemon(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID del Pokémon desde la URL
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID del Pokémon desde la URL
	path := strings.TrimPrefix(r.URL.Path, "/pokemon/")
	if path == "" || strings.Contains(path, "/") {
		http.Error(w, "Ruta inválida", http.StatusBadRequest)
		return
	}
	id := path

	response, err := getPokemonAndAbility(id)
	if err != nil {
		log.Fatalf("Erroooor")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handlerAllPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var id string

	path := strings.TrimPrefix(r.URL.Path, "/all/")
	id = path

	response, err := getAllPokemons(id)
	if err != nil {
		log.Fatalf("Erroooor")
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func getAllPokemons(quantity string) (*ResponseAllPokemonsApi, error) {
	const defaultLimit = "1100"

	var limit string
	if quantity != "" {
		limit = quantity
	} else {
		limit = defaultLimit
	}
	//PRIMER LLAMADO
	//************************************************************************************************************

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/?limit=%s&offset=0", limit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Estado HTTP no esperado: %d", resp.StatusCode)
	}

	pokemon := &ResponseAllPokemons{}
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		log.Fatalf("Error al decodificar JSON: %v", err)
	}
	//************************************************************************************************************
	resApiList := make([]ResponseApi, 0)
	response := ResponseAllPokemonsApi{}
	var total int
	wg := sync.WaitGroup{}
	//list := make(chan []ResponseApi)
	if pokemon.Results != nil {
		for i := 0; i < len(pokemon.Results); i++ {
			wg.Add(1)
			respAllPokemonsAndAbilities, err := getPokemonsAndAbility(pokemon.Results[i].URL, &wg)
			if err != nil {
				log.Fatalf("Error al traer pokemon y habilidad : %v", err)
			}
			resApiList = append(resApiList, respAllPokemonsAndAbilities)
			total = i
		}
		wg.Wait()
		//list <- resApiList

		//<-list
		/*for i := 0; i < len(pokemon.Results); i++ {
			respAllPokemonsAndAbilities, err := getPokemonsAndAbility(pokemon.Results[i].URL)
			if err != nil {
				log.Fatalf("Error al traer pokemon y habilidad : %v", err)
			}
			resApiList = append(resApiList, respAllPokemonsAndAbilities)
			total = i
		}*/
		response.ResponseApi = &resApiList
		response.Total += total + 1
	}
	return &response, nil
}

//func worker(id int, jobs <-chan []Results, results chan<- []Results, wg *sync.WaitGroup) {
//	defer wg.Done()
//	for value := range jobs {
//		fmt.Printf("Worker %d procesando %s\n", id, value)
//		result, err := callAPI(value)
//		if err != nil {
//			results <- fmt.Sprintf("Error al procesar %s: %v", value, err)
//		} else {
//			results <- result
//		}
//	}
//}

func getPokemonsAndAbility(url string, wg *sync.WaitGroup) (ResponseApi, error) {
	defer wg.Done()
	response := ResponseApi{}
	resp, err := http.Get(url)
	if err != nil {
		return ResponseApi{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Estado HTTP no esperado: %d", resp.StatusCode)
	}

	pokemon := &ResponsePokeApi{}
	var abilities []string
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		log.Fatalf("Error al decodificar JSON: %v", err)
	}
	if pokemon.Abilities != nil {
		for i := 0; i < len(pokemon.Abilities); i++ {
			abilities = append(abilities, getAbility(pokemon.Abilities[i].Ability.URL))
		}
		response.Abilityes = abilities
		response.Name = pokemon.Name
		response.Weight = int64(pokemon.Weight)

	}

	return response, nil
}

func getPokemon(id string) (*ResponseApi, error) {

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Estado HTTP no esperado: %d", resp.StatusCode)
	}

	pokemon := &ResponseApi{}
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		log.Fatalf("Error al decodificar JSON: %v", err)
	}

	return pokemon, nil
}

func getPokemonAndAbility(id string) (*ResponseApi, error) {
	response := &ResponseApi{}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Estado HTTP no esperado: %d", resp.StatusCode)
	}

	pokemon := &ResponsePokeApi{}
	var abilities []string
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		log.Fatalf("Error al decodificar JSON: %v", err)
	}
	if pokemon.Abilities != nil {
		for i := 0; i < len(pokemon.Abilities); i++ {
			abilities = append(abilities, getAbility(pokemon.Abilities[i].Ability.URL))
		}
		response.Abilityes = abilities
		response.Name = pokemon.Name
		response.Weight = int64(pokemon.Weight)

	}

	return response, nil
}

func getAbility(url string) string {
	var abilityResponse string

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Estado HTTP no esperado: %d", resp.StatusCode)
	}
	abilityEntity := &AbilityResponse{}
	err = json.NewDecoder(resp.Body).Decode(abilityEntity)
	if err != nil {
		log.Fatalf("Error al decodificar JSON: %v", err)
	}
	abilityResponse = validateResponseAbility(abilityEntity)

	return abilityResponse
}

func validateResponseAbility(ability *AbilityResponse) string {
	var description string
	if ability.FlavorTextEntries == nil {
		return ""
	}
	for i := 0; i < len(ability.FlavorTextEntries); i++ {
		if ability.FlavorTextEntries[i].Language.Name == "es" {
			description = ability.FlavorTextEntries[i].FlavorText
			return description
		}
	}
	return description
}

type ResponseAllPokemons struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous any       `json:"previous"`
	Results  []Results `json:"results"`
}
type Results struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ResponseAllPokemonsApi struct {
	Total       int `json:"total"`
	ResponseApi *[]ResponseApi
}
type ResponseApi struct {
	Name      string
	Weight    int64
	Abilityes []string
}

type AbilityResponse struct {
	FlavorTextEntries []FlavorTextEntries `json:"flavor_text_entries"`
}

type FlavorTextEntries struct {
	FlavorText   string       `json:"flavor_text"`
	Language     Language     `json:"language"`
	VersionGroup VersionGroup `json:"version_group"`
}
type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionGroup struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ResponsePokeApi struct {
	Abilities []Abilities `json:"abilities"`
	Name      string      `json:"name"`
	Weight    int         `json:"weight"`
}

type Abilities struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
