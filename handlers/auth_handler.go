package handlers

import (
	"encoding/json"
	"io"
	"log"
	"log-engine/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	var credentials models.Credentials
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatalln("Error in Reading", err)
	}
	if err := json.Unmarshal(data, &credentials); err != nil {
		log.Fatalln("Error in Creds", err)
	}
	log.Println(credentials.Username)
	log.Println(credentials.Password)
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
