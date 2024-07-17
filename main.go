package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
    initDB()
    defer db.Close()
    r := gin.Default()

    
    personGroup := r.Group("/person")
    {
        personGroup.GET("/:person_id/info", getPersonInfo)
        personGroup.POST("/create", createPerson)
    }

    port := "8080" 
    log.Printf("Server running on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}
