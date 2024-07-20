package service

import (
	"fmt"
	"log/slog"

	"github.com/thetherington/jobber-common/models/notification"
	"github.com/thetherington/jobber-notification/internal/adapters/config"
	"github.com/thetherington/jobber-notification/internal/core/port"
)

const (
	APPICON = "https://i.ibb.co/Kyp2m0t/cover.png"
)

/**
 * NotificationService implements
 */
type NotificationService struct {
	mailer    port.Mailer
	templates map[string]port.TemplateMaker
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(mailer port.Mailer, templates map[string]port.TemplateMaker) *NotificationService {
	return &NotificationService{
		mailer:    mailer,
		templates: templates,
	}
}

func (s *NotificationService) PrintTemplate(template string, locals notification.AuthEmailLocals) {
	html, err := s.templates[template].Render(locals)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(html)
}

func (s *NotificationService) SendAuthEmail(template string, locals notification.AuthEmailLocals) {
	locals.AppIcon = APPICON
	locals.AppLink = config.Config.App.ClientUrl

	htmlBody, err := s.templates[template].Render(locals)
	if err != nil {
		slog.With("error", err).Error("failed to render template", "template", template)
	}

	var subject string

	switch template {
	case "verifyEmail":
		subject = "Welcome to Jobber! Plase Verify Your Email"

	case "forgotPassword":
		subject = "Reset your Jobber Password"

	case "resetPasswordSuccess":
		subject = "Password Reset Successful"

	case "offer":
		subject = "You have received a custom offer from <%= sender %>"
	}

	if err := s.mailer.Send(locals.Username, locals.ReceiverEmail, subject, htmlBody); err != nil {
		slog.With("error", err).Error("failed to send email")
	}
}

func (s *NotificationService) SendOrderEmail(template string, locals notification.OrderEmailLocals) {
	locals.AppIcon = APPICON
	locals.AppLink = config.Config.App.ClientUrl

	htmlBody, err := s.templates[template].Render(locals)
	if err != nil {
		slog.With("error", err).Error("failed to render template", "template", template)
	}

	var (
		subject string
		name    string
		email   string
	)

	switch template {
	case "offer":
		subject = "You have received a custom offer from " + locals.Sender
		name = locals.BuyerUsername
		email = locals.ReceiverEmail

	case "orderPlaced":
		subject = "You've received an order from " + locals.BuyerUsername
		name = locals.SellerUsername
		email = locals.ReceiverEmail

	case "orderReceipt":
		subject = "Hers's your order receipt"
		name = locals.BuyerUsername
		email = locals.ReceiverEmail

	case "orderExtension":
		subject = "You received a delivery extension request from " + locals.SellerUsername
		name = locals.BuyerUsername
		email = locals.ReceiverEmail

	case "orderExtensionApproval":
		subject = locals.Subject
		name = locals.SellerUsername
		email = locals.ReceiverEmail

	case "orderDelivered":
		subject = "Consider it done: Your order is ready for review"
		name = locals.BuyerUsername
		email = locals.ReceiverEmail
	}

	if err := s.mailer.Send(name, email, subject, htmlBody); err != nil {
		slog.With("error", err).Error("failed to send email")
	}
}
