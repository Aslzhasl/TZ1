package model

import (
	"TZ/internal/db"
	"os"
	"testing"
)

func TestPatchPersonByID(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "test_db")

	db.InitPostgres()

	p := Person{
		Name:        "Test",
		Surname:     "User",
		Patronymic:  "Mid",
		Age:         30,
		Gender:      "male",
		Nationality: "US",
	}
	err := InsertPerson(&p)
	if err != nil {
		t.Fatalf("InsertPerson failed: %v", err)
	}

	update := map[string]interface{}{
		"surname": "Updated",
		"age":     35,
	}

	err = PatchPersonByID(int(p.ID), update)
	if err != nil {
		t.Fatalf("PatchPersonByID failed: %v", err)
	}

	got, err := GetPersonByID(int(p.ID))
	if err != nil {
		t.Fatalf("GetPersonByID failed: %v", err)
	}

	if got.Surname != "Updated" || got.Age != 35 {
		t.Errorf("Unexpected updated result: %+v", got)
	}
}
