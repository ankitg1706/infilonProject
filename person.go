package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

func getPersonInfo(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("person_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	person, err := getPersonByID(personID)
	if err != nil {
		log.Println("Error fetching person info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch person info"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func createPerson(c *gin.Context) {
	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := insertPerson(person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
}


