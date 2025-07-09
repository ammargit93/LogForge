package handlers

import (
	"encoding/json"
	"io"
	"log"
	"log-engine/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func credsFilePath() string {
// 	home, _ := os.UserHomeDir()
// 	return home + "/.logcli/creds.json"
// }
// func GetCreds() (string, string) {
// 	path := credsFilePath()
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		log.Println("Error in getting creds", err)
// 	}
// 	var creds Credentials
// 	if err := json.Unmarshal(data, &creds); err != nil {
// 		log.Fatalln("Error Unmarshalling ", err)
// 	}

// 	return creds.Username, creds.Password
// }

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
