package rabbitmq

import (
	rabbitmq "github.com/thetherington/jobber-common/queue"
	"github.com/thetherington/jobber-notification/internal/core/port"
)

type Queue struct {
	service port.NotificationService

	rabbitmq.RabbitQueueAdapter
}

func NewRabbitMQAdapter(url string, svc port.NotificationService) *Queue {
	queue := new(Queue)

	queue.service = svc

	queue.URL = url
	queue.Consumers = make([]rabbitmq.Consumer, 0)

	queue.Connect()
	go queue.Reconnector()

	return queue
}
