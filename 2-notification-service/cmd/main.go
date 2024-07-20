package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/thetherington/jobber-common/logger"
	"github.com/thetherington/jobber-notification/internal/adapters/config"
	"github.com/thetherington/jobber-notification/internal/adapters/gomail"
	"github.com/thetherington/jobber-notification/internal/adapters/handler/http"
	"github.com/thetherington/jobber-notification/internal/adapters/handler/rabbitmq"
	"github.com/thetherington/jobber-notification/internal/adapters/mailtemplates"
	"github.com/thetherington/jobber-notification/internal/core/port"
	"github.com/thetherington/jobber-notification/internal/core/service"
)

const (
	WEB_PORT = 5001
	APP_NAME = "notification-service"
)

var TEMPLATES = [...]string{"verifyEmail", "forgotPassword", "resetPasswordSuccess",
	"offer", "orderPlaced", "orderReceipt", "orderExtension", "orderExtensionApproval", "orderDelivered"}

func main() {
	// Load environment variables
	config := config.Config

	// Set logger
	logger.Set(config.App.Env, APP_NAME, config.App.LogLevel)

	// Create a template maker for each template file
	templateMakers := make(map[string]port.TemplateMaker)

	for _, template := range TEMPLATES {
		tm, err := mailtemplates.NewTemplateMaker(fmt.Sprintf("%s.tmpl.html", template))
		if err != nil {
			slog.With("error", err).Error("Failed to import template", "template", fmt.Sprintf("%s.tmpl.html", template))
			os.Exit(1)
		}

		templateMakers[template] = tm
	}

	// Create a mailer client
	mailer, err := gomail.NewMailClient("smtp.ethereal.email", config.Email.Sender, config.Email.Password)
	if err != nil {
		slog.With("error", err).Error("Failed to create mail client")
	}

	mailer.SetFrom("Jobber App", config.Email.Sender)

	// Create the notification service with dependency injection for mail client and templates
	svc := service.NewNotificationService(mailer, templateMakers)

	// Create the RabbitMQ connection and add the consumers for authentication service and order service
	consumer := rabbitmq.NewRabbitMQAdapter(config.RabbitMQ.Endpoint, svc)
	defer consumer.Close()

	if err := consumer.AddConsumer(consumer.ConsumeAuthEmailMessages); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeAuthEmailMessages")
		os.Exit(1)
	}

	if err := consumer.AddConsumer(consumer.ConsumeOrderEmailMessages); err != nil {
		slog.With("error", err).Error("Failed to Add ConsumeOrderEmailMessages")
		os.Exit(1)
	}

	// setup the http server to listen for Ping health checks
	router := http.NewRouter()

	addr := fmt.Sprintf(":%d", WEB_PORT)

	slog.Info("Starting HTTP Server", "address", addr)
	if err := http.Serve(addr, router); err != nil {
		slog.With("error", err).Error("Failed to start HTTP Server")
	}
}
