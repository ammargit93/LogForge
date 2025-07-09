package handlers

import (
	"fmt"
	"log"
	"log-engine/parquet"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Message string `json:"message"`
}

func HandleQuery(c *gin.Context) {
	var query Query
	username, _ := parquet.GetCreds()

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dataPath := extractParquetPath(query.Message)
	arr := strings.Split(dataPath, "/")[0]
	log.Println(arr)
	log.Println(username)
	if username != arr {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad username"})
		return
	}
	log.Println("Received SQL Query:", query.Message)

	if dataPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not extract path from SQL query"})
		return
	}

	// Step 2: Build DuckDB query with proper S3 configuration
	finalQuery := buildDuckDBQuery(query.Message, dataPath)

	// Step 3: Run DuckDB CLI with built query
	output := runDuckDBQuery(finalQuery)

	// Step 4: Return output to client

	log.Println(output)
	c.JSON(http.StatusOK, gin.H{
		"message": "Query executed successfully",
		"output":  output,
	})
}

func extractParquetPath(sql string) string {
	re := regexp.MustCompile(`(?i)FROM\s+'([^']+)'`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func buildDuckDBQuery(originalQuery, dataPath string) string {
	config := `
SET s3_region='us-east-1';
SET s3_access_key_id='minioadmin';
SET s3_secret_access_key='minioadmin';
SET s3_endpoint='localhost:9000';
SET s3_url_style='path';
SET s3_use_ssl=false;
`
	s3Path := fmt.Sprintf("s3://logs/%s", dataPath)
	modifiedQuery := strings.Replace(originalQuery,
		fmt.Sprintf("FROM '%s'", dataPath),
		fmt.Sprintf("FROM read_parquet('%s')", s3Path), 1)

	return config + modifiedQuery
}

func runDuckDBQuery(query string) string {
	cmd := exec.Command("C:\\Users\\Ammar1\\go\\log-engine\\bin\\duckdb.exe", "-c", query)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ DuckDB error: %v\n", err)
		log.Printf("🔴 Output: %s\n", string(output))
		return fmt.Sprintf("DuckDB Error: %v\nOutput: %s", err, string(output))
	}
	return string(output)
}
