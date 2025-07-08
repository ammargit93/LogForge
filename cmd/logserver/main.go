package main

import (
	"log"
	"log-engine/handlers"

	"github.com/gin-gonic/gin"
)

// Finish the query route.
// login and logout come next
// ai engine (simple langchain)
// caching frequently accessed files.(redis, emberdb)
// ppt and docs (optional)
// Fix minor bugs

func main() {
	router := gin.Default()
	router.POST("/logs", handlers.HandleLog)
	// query route
	log.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}
