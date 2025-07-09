package main

import (
	"log"
	"log-engine/handlers"

	"github.com/gin-gonic/gin"
)

// login and logout come next
// ai engine (simple langchain)
// caching frequently accessed files.(redis, emberdb)
// ppt and docs (optional)
// Fix minor bugs

func main() {
	router := gin.Default()
	router.POST("/logs", handlers.HandleLog)
	router.POST("/query", handlers.HandleQuery)
	router.POST("/login", handlers.HandleLogin)

	log.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}
