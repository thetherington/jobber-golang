package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thetherington/jobber-common/models/review"
	pbReview "github.com/thetherington/jobber-common/protogen/go/review"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/protobuf/proto"
)

func (q *Queue) ConsumeGigDirectMessage() error {
	var (
		exchangeName = "jobber-update-gig"
		routingKey   = "update-gig"
		queueName    = "gig-update-queue"
	)

	ch, err := q.Connection.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
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

	err = ch.QueueBind(queue.Name, routingKey, exchangeName, false, amqp091.Table{})
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
				slog.With("error", err).Error("ConsumeGigDirectMessage: failed to unmarshal protobuf message")
				continue
			}

			if pbMsg.Action == pbReview.ReviewType_GigReview {
				err := q.gig.UpdateGigReview(&review.ReviewMessageDetails{
					GigId:      pbMsg.GigId,
					ReviewerId: pbMsg.ReviewerId,
					SellerId:   pbMsg.SellerId,
					Review:     pbMsg.Review,
					Rating:     pbMsg.Rating,
					OrderId:    pbMsg.OrderId,
					Type:       pbMsg.Action.String(),
					CreatedAt:  *utils.ToTime(pbMsg.CreatedAt),
				})
				if err != nil {
					slog.With("error", err).Error("failed to update gig review")
				}
			}

			msg.Ack(false)
		}

		slog.Debug("ConsumeGigDirectMessage shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}

func (q *Queue) ConsumeSeedDirectMessage() error {
	var (
		exchangeName = "jobber-seed-gig"
		routingKey   = "receive-sellers"
		queueName    = "seed-gig-queue"
	)

	ch, err := q.Connection.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
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

	err = ch.QueueBind(queue.Name, routingKey, exchangeName, false, amqp091.Table{})
	if err != nil {
		return err
	}

	messages, err := ch.Consume(queue.Name, "", false, false, false, false, amqp091.Table{})
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			var pbMsg pb.SellersResponse

			if err := proto.Unmarshal(msg.Body, &pbMsg); err != nil {
				slog.With("error", err).Error("ConsumeSeedDirectMessage: failed to unmarshal protobuf message")
				continue
			}

			// use any so that the ports interface doesn't have to hard implement the protobuf message
			sellers := make([]any, 0)
			for _, s := range pbMsg.Sellers {
				sellers = append(sellers, s)
			}

			q.gig.SeedData(context.Background(), sellers)

			msg.Ack(false)
		}

		slog.Debug("ConsumeSeedDirectMessage shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}
