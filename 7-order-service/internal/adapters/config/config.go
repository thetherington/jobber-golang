package config

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type (
	Container struct {
		App        *App
		Tokens     *Tokens
		RabbitMQ   *RabbitMQ
		Stripe     *Stripe
		DB         *DB
		Redis      *Redis
		Cloudinary *Cloudinary
		Elastic    *Elastic
	}

	App struct {
		APM        bool
		Env        string
		ClientUrl  string
		GatewayUrl string
		LogLevel   slog.Level
	}

	Tokens struct {
		GATEWAY string
		JWT     string
	}

	RabbitMQ struct {
		Endpoint string
	}

	Stripe struct {
		Key string
	}

	DB struct {
		URI  string
		Name string
	}

	Redis struct {
		Host string
	}

	Cloudinary struct {
		Name      string
		ApiKey    string
		ApiSecret string
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
		Env:        os.Getenv("APP_ENV"),
		ClientUrl:  os.Getenv("CLIENT_URL"),
		GatewayUrl: os.Getenv("API_GATEAWY_URL"),
		APM:        false,
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

	tokens := &Tokens{
		GATEWAY: os.Getenv("GATEWAY_JWT_TOKEN"),
		JWT:     os.Getenv("JWT_TOKEN"),
	}

	rabbitmq := &RabbitMQ{
		Endpoint: os.Getenv("RABBITMQ_ENDPOINT"),
	}

	stripe := &Stripe{
		Key: os.Getenv("STRIPE_API_KEY"),
	}

	db := &DB{
		URI:  os.Getenv("DB_URI"),
		Name: os.Getenv("DB_NAME"),
	}

	redis := &Redis{
		Host: os.Getenv("REDIS_HOST"),
	}

	cloud := &Cloudinary{
		Name:      os.Getenv("CLOUD_NAME"),
		ApiKey:    os.Getenv("CLOUD_API_KEY"),
		ApiSecret: os.Getenv("CLOUD_API_SECRET"),
	}

	elastic := &Elastic{
		SearchUrl:      os.Getenv("ELASTIC_SEARCH_URL"),
		ApmUrl:         os.Getenv("ELASTIC_APM_SERVER_URL"),
		ApmSecretToken: os.Getenv("ELASTIC_APM_SECRET_TOKEN"),
	}

	Config = &Container{
		app,
		tokens,
		rabbitmq,
		stripe,
		db,
		redis,
		cloud,
		elastic,
	}
}
