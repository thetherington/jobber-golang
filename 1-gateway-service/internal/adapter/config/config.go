package config

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type (
	Container struct {
		App      *App
		Tokens   *Tokens
		Secrets  *Secrets
		Services *Services
		Redis    *Redis
		Elastic  *Elastic
	}

	App struct {
		APM       bool
		Env       string
		ClientUrl string
		LogLevel  slog.Level
	}

	Tokens struct {
		GATEWAY string
		JWT     string
	}

	Secrets struct {
		KEY_ONE string
		KEY_TWO string
	}

	Services struct {
		Auth    string
		Users   string
		Gig     string
		Message string
		Order   string
		Review  string
	}

	Redis struct {
		Host string
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
		LogLevel:  slog.LevelInfo,
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

	secrets := &Secrets{
		KEY_ONE: os.Getenv("SECRET_KEY_ONE"),
		KEY_TWO: os.Getenv("SECRET_KEY_TWO"),
	}

	services := &Services{
		Auth:    os.Getenv("AUTH_ENDPOINT"),
		Users:   os.Getenv("USERS_ENDPOINT"),
		Gig:     os.Getenv("GIG_ENDPOINT"),
		Message: os.Getenv("MESSAGE_ENDPOINT"),
		Order:   os.Getenv("ORDER_ENDPOINT"),
		Review:  os.Getenv("REVIEW_ENDPOINT"),
	}

	redis := &Redis{
		Host: os.Getenv("REDIS_HOST"),
	}

	elastic := &Elastic{
		SearchUrl:      os.Getenv("ELASTIC_SEARCH_URL"),
		ApmUrl:         os.Getenv("ELASTIC_APM_SERVER_URL"),
		ApmSecretToken: os.Getenv("ELASTIC_APM_SECRET_TOKEN"),
	}

	Config = &Container{
		app,
		tokens,
		secrets,
		services,
		redis,
		elastic,
	}
}
