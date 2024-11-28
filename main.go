
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {

	fmt.Println("Servidor escuchando en http://localhost:8080")
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
	var abilieites []string
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		log.Fatalf("Error al decodificar JSON: %v", err)
	}
	if pokemon.Abilities != nil {
		for i := 0; i < len(pokemon.Abilities); i++ {
			abilieites = append(abilieites, getAbility(pokemon.Abilities[i].Ability.URL))
		}
		response.Abilityes = abilieites
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
