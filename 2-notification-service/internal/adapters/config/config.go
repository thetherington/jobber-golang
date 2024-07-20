package config

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type (
	Container struct {
		App      *App
		RabbitMQ *RabbitMQ
		Email    *Email
		Elastic  *Elastic
	}

	App struct {
		APM       bool
		Env       string
		ClientUrl string
		LogLevel  slog.Level
	}

	RabbitMQ struct {
		Endpoint string
	}

	Email struct {
		Sender   string
		Password string
	}

	Elastic struct {
		SearchUrl      string
		ApmUrl         string
		ApmSecretToken string
	}
)

var Config *Container

// New creates a new container instance
func init() {
	app := &App{
		Env:       os.Getenv("APP_ENV"),
		ClientUrl: os.Getenv("CLIENT_URL"),
		APM:       false,
	}

	if apm := os.Getenv("ENABLE_APM"); apm == "1" {
		app.APM = true
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		switch logLevel {
		case "debug":
			app.LogLevel = slog.LevelDebug
		case "info":
			app.LogLevel = slog.LevelInfo
		}
	}

	rabbitmq := &RabbitMQ{
		Endpoint: os.Getenv("RABBITMQ_ENDPOINT"),
	}

	email := &Email{
		Sender:   os.Getenv("SENDER_EMAIL"),
		Password: os.Getenv("SENDER_EMAIL_PASSWORD"),
	}

	elastic := &Elastic{
		SearchUrl:      os.Getenv("ELASTIC_SEARCH_URL"),
		ApmUrl:         os.Getenv("ELASTIC_APM_SERVER_URL"),
		ApmSecretToken: os.Getenv("ELASTIC_APM_SECRET_TOKEN"),
	}

	Config = &Container{
		app,
		rabbitmq,
		email,
		elastic,
	}
}
