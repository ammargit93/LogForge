package handlers

import (
	"log-engine/models"
	"log-engine/parquet"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleLog(c *gin.Context) {
	var logEntry models.LogEntry
	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username, _ := parquet.GetCreds()
	models.Mu.Lock()
	models.BufferQueue = append(models.BufferQueue, logEntry)
	if len(models.BufferQueue) >= models.N {
		bufferCopy := make([]models.LogEntry, len(models.BufferQueue))
		copy(bufferCopy, models.BufferQueue)
		go parquet.WriteToParquet(bufferCopy, logEntry.Service, username)
		models.BufferQueue = []models.LogEntry{}
	}
	models.Mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
