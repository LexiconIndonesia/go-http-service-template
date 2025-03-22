package main

import (
	"fmt"
	"os"
	"strconv"
)

func loadEnvString(key string, result *string) {
	s, ok := os.LookupEnv(key)

	if !ok {
		return
	}
	*result = s
}

func loadEnvUint(key string, result *uint) {
	s, ok := os.LookupEnv(key)

	if !ok {
		return
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	*result = uint(n)
}

/* Configuration */

/* PgSQL Configuration */
type pgSqlConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	SslMode  string `json:"ssl_mode"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (p pgSqlConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=%s", p.Host, p.Port, p.User, p.Password, p.Database, p.SslMode)
}

func defaultPgSql() pgSqlConfig {
	return pgSqlConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "database",
		User:     "",
		Password: "",
		SslMode:  "disable",
	}
}

func (p *pgSqlConfig) loadFromEnv() {
	loadEnvString("POSTGRES_HOST", &p.Host)
	loadEnvUint("POSTGRES_PORT", &p.Port)
	loadEnvString("POSTGRES_DB_NAME", &p.Database)
	loadEnvString("POSTGRES_SSLMODE", &p.SslMode)
	loadEnvString("POSTGRES_USERNAME", &p.User)
	loadEnvString("POSTGRES_PASSWORD", &p.Password)
}

/* Listen Configuration */

type listenConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

func (l listenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func defaultListenConfig() listenConfig {
	return listenConfig{
		Host: "127.0.0.1",
		Port: 8080,
	}
}

func (l *listenConfig) loadFromEnv() {
	loadEnvString("LISTEN_HOST", &l.Host)
	loadEnvUint("LISTEN_PORT", &l.Port)
}

type hostConfig struct {
	// Add any necessary fields for host configuration
}

func (h *hostConfig) loadFromEnv() {
	// Implement loading from environment variables
}

type natsConfig struct {
	URL      string
	Username string
	Password string
}

func (c *natsConfig) loadFromEnv() {
	c.URL = getEnv("NATS_URL", "nats://localhost:4222")
	c.Username = getEnv("NATS_USERNAME", "")
	c.Password = getEnv("NATS_PASSWORD", "")
}

type config struct {
	Host          hostConfig
	Listen        listenConfig
	PgSql         pgSqlConfig
	BackendApiKey string
	ServerSalt    string
	Nats          natsConfig
}

func (c *config) loadFromEnv() {
	c.Host.loadFromEnv()
	c.Listen.loadFromEnv()
	c.PgSql.loadFromEnv()
	c.BackendApiKey = getEnv("BACKEND_API_KEY", "")
	c.ServerSalt = getEnv("SERVER_SALT", "")
	c.Nats.loadFromEnv()
}

func defaultConfig() config {
	return config{
		Host: hostConfig{
			// ... existing code ...
		},
		Listen: listenConfig{
			// ... existing code ...
		},
		PgSql: pgSqlConfig{
			// ... existing code ...
		},
		BackendApiKey: "",
		ServerSalt:    "",
		Nats: natsConfig{
			URL:      "nats://localhost:4222",
			Username: "",
			Password: "",
		},
	}
}
