package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thetherington/jobber-common/models/review"
	pbReview "github.com/thetherington/jobber-common/protogen/go/review"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/protobuf/proto"
)

func (q *Queue) ConsumerReviewFanoutMessages() error {
	var (
		exchangeName = "jobber-review"
		queueName    = "order-review-queue"
	)

	ch, err := q.Connection.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		amqp091.Table{},
	)
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, amqp091.Table{})
	if err != nil {
		return err
	}

	err = ch.QueueBind(queue.Name, "", exchangeName, false, amqp091.Table{})
	if err != nil {
		return err
	}

	messages, err := ch.Consume(queue.Name, "", false, false, false, false, amqp091.Table{})
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			var pbMsg pbReview.ReviewMessageDetails

			if err := proto.Unmarshal(msg.Body, &pbMsg); err != nil {
				slog.With("error", err).Error("ConsumerReviewFanoutMessages: failed to unmarshal protobuf message")
				continue
			}

			if pbMsg.Action == pbReview.ReviewType_BuyerReview || pbMsg.Action == pbReview.ReviewType_SellerReview {
				resp, userTo, err := q.order.UpdateOrderReview(context.Background(), &review.ReviewMessageDetails{
					GigId:      pbMsg.GigId,
					ReviewerId: pbMsg.ReviewerId,
					SellerId:   pbMsg.SellerId,
					Review:     pbMsg.Review,
					Rating:     pbMsg.Rating,
					OrderId:    pbMsg.OrderId,
					Type:       pbMsg.Action.String(),
					CreatedAt:  *utils.ToTime(pbMsg.CreatedAt),
				})
				if err == nil {
					q.grpc.NotifyUpdateOrder(resp.Order, userTo, resp.Message)
				}
			}

			msg.Ack(false)
		}

		slog.Debug("ConsumerReviewFanoutMessages shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}
