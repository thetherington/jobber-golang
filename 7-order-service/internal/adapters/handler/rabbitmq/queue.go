package rabbitmq

import (
	rabbitmq "github.com/thetherington/jobber-common/queue"
	"github.com/thetherington/jobber-order/internal/core/port"
)

type Queue struct {
	order port.OrderService
	grpc  port.GrpcInterface

	rabbitmq.RabbitQueueAdapter
}

func NewRabbitMQAdapter(url string, order port.OrderService, grpc port.GrpcInterface) *Queue {
	queue := new(Queue)

	queue.order = order
	queue.grpc = grpc

	queue.URL = url
	queue.Consumers = make([]rabbitmq.Consumer, 0)

	queue.Connect()
	go queue.Reconnector()

	return queue
}
