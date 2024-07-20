package rabbitmq

import (
	rabbitmq "github.com/thetherington/jobber-common/queue"
	"github.com/thetherington/jobber-gig/internal/core/port"
)

type Queue struct {
	gig port.GigService

	rabbitmq.RabbitQueueAdapter
}

func NewRabbitMQAdapter(url string, gig port.GigService) *Queue {
	queue := new(Queue)

	queue.gig = gig

	queue.URL = url
	queue.Consumers = make([]rabbitmq.Consumer, 0)

	queue.Connect()
	go queue.Reconnector()

	return queue
}
