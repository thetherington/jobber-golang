package port

import (
	"github.com/thetherington/jobber-common/models/notification"
)

type NotificationService interface {
	PrintTemplate(template string, locals notification.AuthEmailLocals)
	SendAuthEmail(template string, locals notification.AuthEmailLocals)
	SendOrderEmail(template string, locals notification.OrderEmailLocals)
}

type Mailer interface {
	Send(toName string, toEmail string, subject string, msgBody string) error
}

type TemplateMaker interface {
	Render(data any) (string, error)
}
