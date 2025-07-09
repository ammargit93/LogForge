package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func credsFilePath() string {
	home, _ := os.UserHomeDir()
	return home + "/.logcli/creds.json"
}

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
		Usage: "Send and query logs on a remote server",
		Commands: []*cli.Command{
			{
				Name:  "send",
				Usage: "Send a log entry",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "timestamp", Usage: "Log timestamp (RFC3339)", Required: true},
					&cli.StringFlag{Name: "level", Usage: "Log level", Required: true},
					&cli.StringFlag{Name: "message", Usage: "Log message", Required: true},
					&cli.StringFlag{Name: "service", Usage: "Service name", Required: true},
					&cli.StringFlag{Name: "host", Usage: "Host identifier", Required: true},
					&cli.StringFlag{Name: "server", Usage: "Server URL", Value: "http://localhost:8080/logs"},
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
						return err
					}
					fmt.Println("✅ Log sent successfully")
					return nil
				},
			},
			{
				Name:  "query",
				Usage: "Query logs using a filter expression",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "q", Usage: `Query string (e.g., level="ERROR" AND service="auth-service")`, Required: true},
					&cli.StringFlag{Name: "server", Usage: "Query endpoint", Value: "http://localhost:8080/query"},
				},
				Action: func(c *cli.Context) error {
					query := map[string]string{"message": c.String("q")}

					jsonData, err := json.Marshal(query)
					if err != nil {
						return fmt.Errorf("failed to encode query: %w", err)
					}

					resp, err := http.Post(c.String("server"), "application/json", bytes.NewBuffer(jsonData))
					if err != nil {
						return fmt.Errorf("request failed: %w", err)
					}
					defer resp.Body.Close()

					body, _ := io.ReadAll(resp.Body)
					if resp.StatusCode >= 400 {
						return fmt.Errorf("server error %d: %s", resp.StatusCode, string(body))
					}

					// Parse and print the output field from the server response
					var result map[string]interface{}
					if err := json.Unmarshal(body, &result); err != nil {
						return fmt.Errorf("failed to parse response: %w", err)
					}
					if output, ok := result["output"]; ok {
						// fmt.Println("────────────────────────────────────")
						fmt.Println(output)
					} else {
						fmt.Println("❌ No 'output' field found in server response:")
						fmt.Println(string(body))
					}

					return nil
				},
			},
			{
				Name:  "login",
				Usage: "Log in and store credentials",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "username", Usage: "Username", Required: true},
					&cli.StringFlag{Name: "password", Usage: "Password", Required: true},
				},
				Action: func(c *cli.Context) error {
					creds := Credentials{
						Username: c.String("username"),
						Password: c.String("password"),
					}
					path := credsFilePath()
					os.MkdirAll(filepath.Dir(path), 0700)
					data, _ := json.Marshal(creds)
					err := os.WriteFile(path, data, 0600)
					if err != nil {
						return fmt.Errorf("failed to save credentials: %w", err)
					}
					resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(data))
					if err != nil {
						return err
					}
					defer resp.Body.Close()

					fmt.Println("✅ Logged in successfully.")
					return nil
				},
			},
			{
				Name:  "logout",
				Usage: "Clear stored credentials",
				Action: func(c *cli.Context) error {
					path := credsFilePath()
					if err := os.Remove(path); err != nil {
						return fmt.Errorf("failed to delete credentials: %w", err)
					}
					fmt.Println("✅ Logged out successfully.")
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
