package handler

import (
	"TZ/internal/db"
	"TZ/internal/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCreatePerson(t *testing.T) {

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "yourpassword")
	os.Setenv("DB_NAME", "people_test")

	db.InitPostgres()
	gin.SetMode(gin.TestMode)

	body := `{
		"name": "Dmitriy",
		"surname": "Ushakov",
		"patronymic": "Vasilevich"
	}`

	req := httptest.NewRequest(http.MethodPost, "/people", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	CreatePerson(c)
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", w.Code)
	}

	var person model.Person
	if err := json.Unmarshal(w.Body.Bytes(), &person); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if person.Age == 0 || person.Gender == "" || person.Nationality == "" {
		t.Errorf("Enrichment failed: %+v", person)
	}
}
