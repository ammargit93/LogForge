package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	Host      string    `json:"host"`
}

func main() {
	app := &cli.App{
		Name:  "logcli",
		Usage: "Send logs to a remote server",
		Commands: []*cli.Command{
			{
				Name:  "send",
				Usage: "Send a log entry",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "timestamp",
						Usage:    "Log timestamp in RFC3339 format",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "level",
						Usage:    "Log level (INFO, WARN, ERROR)",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "message",
						Usage:    "Log message content",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "service",
						Usage:    "Service name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "host",
						Usage:    "Host identifier",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "server",
						Usage: "Log server endpoint",
						Value: "http://localhost:8080/logs",
					},
				},
				Action: func(c *cli.Context) error {
					timestamp, err := time.Parse(time.RFC3339, c.String("timestamp"))
					if err != nil {
						return fmt.Errorf("invalid timestamp format: %w", err)
					}

					entry := LogEntry{
						Timestamp: timestamp,
						Level:     c.String("level"),
						Message:   c.String("message"),
						Service:   c.String("service"),
						Host:      c.String("host"),
					}

					if err := sendLog(c.String("server"), entry); err != nil {
						return fmt.Errorf("failed to send log: %w", err)
					}

					fmt.Println("Log sent successfully")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func sendLog(serverURL string, entry LogEntry) error {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}
