# 🧠 Log Query Engine with DuckDB + MinIO

A lightweight log processing system using **DuckDB**, **MinIO**, and a **Go-based CLI client**.

> ✅ Send logs → 🗃️ Stored as Parquet in MinIO → 🔎 Query with DuckDB SQL

---

## 📌 Overview

This project provides:
- A **Go backend server** that receives logs and SQL queries
- A **CLI tool (`logcli`)** to interact with the server
- Logs are stored in **MinIO** as Parquet files
- Queries are executed using **DuckDB** with S3 integration

---

## 🧱 Tech Stack

- **Go (Golang)** – Backend and CLI
- **Gin** – HTTP server
- **DuckDB** – Embedded analytics database
- **MinIO** – S3-compatible object storage
- **Parquet** – Columnar file format
- **urfave/cli** – CLI interface

---

## 💻 Components

### 🔗 Server (API)
Runs on: `http://localhost:8080`  
Endpoints:
- `POST /logs` – Accepts log entries
- `POST /query` – Accepts SQL queries

### 🖥️ Client: CLI Tool (`logcli`)
- Command-line interface to interact with the server
- Sends logs and queries
- Written in Go (cross-platform)

---

## ⚙️ Setup

### 1. ✅ Requirements

- Go 1.18+
- [DuckDB CLI](https://duckdb.org/docs/installation/cli.html)
- MinIO (running on `http://localhost:9000`)
- MinIO credentials:
  - Access Key: `minioadmin`
  - Secret Key: `minioadmin`
 
  
### ▶️ Install MinIO Server

Download the MinIO executable for your platform:

- **Windows (64-bit):**  
  `https://dl.min.io/server/minio/release/windows-amd64/minio.exe` :contentReference[oaicite:4]{index=4}
- **Linux/macOS:**  
  See [MinIO download page](https://min.io) for latest binaries.
Place `minio.exe` inside your project’s `bin/` folder.

---

### ▶️ Start MinIO

```bash
minio server ./data --address ":9000" --console-address ":9001"
```
To build the client exe
```bash
go build -o logcli client/cmd/main.go
```

To run The server

```bash
go run /cmd/logserver/main.go
```
### Note 
Note: This project is merely an MVP.
