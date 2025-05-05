package handler

import (
	"TZ/internal/model"
	"TZ/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetallPerson(c *gin.Context) {
	name := c.Query("name")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	people, err := model.GetPeople(name, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch people"})
		return
	}

	c.JSON(http.StatusOK, people)
}

func CreatePerson(c *gin.Context) {
	var p model.Person
	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("[WARN] Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[INFO] Received person: %+v", p)

	if err := service.EnrichPerson(&p); err != nil {
		log.Printf("[ERROR] Failed to enrich person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enrich person"})
		return
	}
	log.Printf("[INFO] Enriched person: %+v", p)

	if err := model.InsertPerson(&p); err != nil {
		log.Printf("[ERROR] Failed to save person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save person"})
		return
	}

	log.Printf("[INFO] Person saved: %+v", p)
	c.JSON(http.StatusCreated, p)
}

func GetPersonByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	person, err := model.GetPersonByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}
func DeletePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = model.DeletePersonByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "person deleted successfully"})
}

func UpdatePerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var p model.Person
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	p.ID = int(int64(id))

	err = model.UpdatePerson(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "person updated successfully"})
}
func PatchPerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[WARN] Invalid ID: %v", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[WARN] Invalid JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if len(input) == 0 {
		log.Printf("[INFO] Empty patch payload for ID=%d", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty payload"})
		return
	}

	log.Printf("[INFO] Patching person ID=%d with fields: %+v", id, input)

	if err := model.PatchPersonByID(id, input); err != nil {
		log.Printf("[ERROR] Failed to patch person ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update person"})
		return
	}

	log.Printf("[INFO] Successfully patched person ID=%d", id)
	c.JSON(http.StatusOK, gin.H{"message": "person patched successfully"})
}
