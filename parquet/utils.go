package parquet

import (
	"encoding/json"
	"log"
	"log-engine/models"
	"os"
)

func credsFilePath() string {
	home, _ := os.UserHomeDir()
	return home + "/.logcli/creds.json"
}
func GetCreds() (string, string) {
	path := credsFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error in getting creds", err)
	}
	var creds models.Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		log.Fatalln("Error Unmarshalling ", err)
	}

	return creds.Username, creds.Password
}
