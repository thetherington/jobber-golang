package rabbitmq

import (
	rabbitmq "github.com/thetherington/jobber-common/queue"
	"github.com/thetherington/jobber-users/internal/core/port"
)

type Queue struct {
	seller port.SellerService
	buyer  port.BuyerService

	rabbitmq.RabbitQueueAdapter
}

func NewRabbitMQAdapter(url string, seller port.SellerService, buyer port.BuyerService) *Queue {
	queue := new(Queue)

	queue.seller = seller
	queue.buyer = buyer

	queue.URL = url
	queue.Consumers = make([]rabbitmq.Consumer, 0)

	queue.Connect()
	go queue.Reconnector()

	return queue
}
