package models

import "sync"

type LogEntry struct {
	Timestamp string `json:"timestamp" binding:"required" parquet:"name=timestamp, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Level     string `json:"level" binding:"required" parquet:"name=level, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Message   string `json:"message" binding:"required" parquet:"name=message, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
	Service   string `json:"service" binding:"required" parquet:"name=service, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Host      string `json:"host" binding:"required" parquet:"name=host, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

var BufferQueue []LogEntry
var Mu sync.Mutex

const N = 5
