package model

import (
	"TZ/internal/db"
	"database/sql"
	"fmt"
	"log"
)

func InsertPerson(p *Person) error {
	log.Printf("[DEBUG] Inserting person: %+v", p)
	_, err := db.DB.Exec(`
		INSERT INTO people (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality)

	if err != nil {
		log.Printf("[ERROR] Insert failed: %v", err)
	}
	return err
}

func GetPeople(name string, limit, offset int) ([]Person, error) {
	var rows *sql.Rows
	var err error

	if name != "" {
		rows, err = db.DB.Query(`
			SELECT id, name, surname, patronymic, age, gender, nationality
			FROM people
			WHERE name ILIKE $1
			ORDER BY id
			LIMIT $2 OFFSET $3
		`, "%"+name+"%", limit, offset)
	} else {
		rows, err = db.DB.Query(`
			SELECT id, name, surname, patronymic, age, gender, nationality
			FROM people
			ORDER BY id
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
			return nil, err
		}
		people = append(people, p)
	}

	return people, nil
}

func GetPersonByID(id int) (*Person, error) {
	row := db.DB.QueryRow("SELECT id, name, surname, patronymic, age, gender, nationality FROM people WHERE id = $1", id)

	var p Person
	err := row.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func DeletePersonByID(id int) error {
	_, err := db.DB.Exec("DELETE FROM people WHERE id = $1", id)
	return err
}

func UpdatePerson(p *Person) error {
	query := `
		UPDATE people 
		SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 
		WHERE id = $7
	`
	_, err := db.DB.Exec(query, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality, p.ID)
	return err
}
func PatchPersonByID(id int, fields map[string]interface{}) error {
	log.Printf("[DEBUG] Patch request for ID=%d: %+v", id, fields)

	allowed := map[string]bool{
		"name": true, "surname": true, "patronymic": true,
		"age": true, "gender": true, "nationality": true,
	}

	query := "UPDATE people SET "
	args := []interface{}{}
	i := 1

	for k, v := range fields {
		if !allowed[k] {
			log.Printf("[WARN] Ignoring invalid field: %s", k)
			continue
		}
		query += fmt.Sprintf("%s = $%d, ", k, i)
		args = append(args, v)
		i++
	}

	if len(args) == 0 {
		log.Printf("[WARN] No valid fields provided for patch")
		return fmt.Errorf("no valid fields provided")
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	log.Printf("[DEBUG] Executing SQL: %s with args: %+v", query, args)
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Printf("[ERROR] Patch SQL failed for ID=%d: %v", id, err)
	}
	return err
}
