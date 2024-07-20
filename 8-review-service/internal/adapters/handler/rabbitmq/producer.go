package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q *Queue) PublishFanoutMessage(exchangeName string, data []byte) error {
	err := q.Channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return err
	}

	err = q.Channel.Publish(exchangeName, "", true, false, amqp.Publishing{
		ContentType:   "application/x-protobuf",
		DeliveryMode:  amqp.Persistent,
		CorrelationId: fmt.Sprintf("users_%d", 69),
		Body:          data,
	})
	if err != nil {
		return err
	}

	return nil
}
