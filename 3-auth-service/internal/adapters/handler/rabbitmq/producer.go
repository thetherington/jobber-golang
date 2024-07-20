package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q *Queue) PublishDirectMessage(exchangeName string, routingKey string, data []byte) error {
	err := q.Channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return err
	}

	err = q.Channel.Publish(exchangeName, routingKey, true, false, amqp.Publishing{
		ContentType:   "application/x-protobuf",
		DeliveryMode:  amqp.Persistent,
		CorrelationId: fmt.Sprintf("user_created_%d", 69),
		Body:          data,
	})
	if err != nil {
		return err
	}

	return nil
}
