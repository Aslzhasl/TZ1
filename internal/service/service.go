package service

import (
	"TZ/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func EnrichPerson(p *model.Person) error {
	name := p.Name
	log.Printf("[INFO] Enriching person: %s", name)

	age, err := fetchAge(name)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch age: %v", err)
		return fmt.Errorf("failed to fetch age: %v", err)
	}
	log.Printf("[DEBUG] Age for %s: %d", name, age)

	gender, err := fetchGender(name)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch gender: %v", err)
		return fmt.Errorf("failed to fetch gender: %v", err)
	}
	log.Printf("[DEBUG] Gender for %s: %s", name, gender)

	nationality, err := fetchNationality(name)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch nationality: %v", err)
		return fmt.Errorf("failed to fetch nationality: %v", err)
	}
	log.Printf("[DEBUG] Nationality for %s: %s", name, nationality)

	p.Age = age
	p.Gender = gender
	p.Nationality = nationality

	log.Printf("[INFO] Enrichment complete for %s", name)
	return nil
}

func fetchAge(name string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Age, nil
}

func fetchGender(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Gender, nil
}

func fetchNationality(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Country) > 0 {
		return data.Country[0].CountryID, nil
	}
	return "", nil
}
