package rabbitmq

import (
	rabbitmq "github.com/thetherington/jobber-common/queue"
)

type Queue struct {
	rabbitmq.RabbitQueueAdapter
}

func NewRabbitMQAdapter(url string) *Queue {
	queue := new(Queue)

	queue.URL = url
	queue.Consumers = make([]rabbitmq.Consumer, 0)

	queue.Connect()
	go queue.Reconnector()

	return queue
}
