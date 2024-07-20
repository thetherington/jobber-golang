package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thetherington/jobber-common/models/users"
	pbEmail "github.com/thetherington/jobber-common/protogen/go/notification"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/protobuf/proto"
)

func (q *Queue) ConsumeBuyerDirectMessage() error {
	var (
		exchangeName = "jobber-buyer-update"
		routingKey   = "user-buyer"
		queueName    = "user-buyer-queue"
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
			var pbMsg pb.BuyerPayload

			if err := proto.Unmarshal(msg.Body, &pbMsg); err != nil {
				slog.With("error", err).Error("Buyer Consumer: failed to unmarshal protobuf message")
				continue
			}

			// Message came from auth service to create a new buyer (user signup / seed)
			if pbMsg.Action.String() == pb.Action_AUTH.String() {
				buyer := &users.Buyer{
					Username:       *pbMsg.Username,
					Email:          *pbMsg.Email,
					ProfilePicture: *pbMsg.ProfilePicture,
					Country:        *pbMsg.Country,
					IsSeller:       false,
					PurchasedGigs:  make([]string, 0),
					CreatedAt:      utils.ToTime(pbMsg.CreatedAt),
					UpdatedAt:      utils.ToTime(pbMsg.UpdatedAt),
				}

				go q.buyer.CreateBuyer(buyer)
			}

			// Message came from Order service for to create/cancel a gig purchase
			if pbMsg.Action.String() == pb.Action_PURCHASED_GIG.String() ||
				pbMsg.Action.String() == pb.Action_CANCELLED_GIG.String() {

				if len(pbMsg.PurchasedGigs) > 0 {
					go q.buyer.UpdateBuyerPurchasedGigs(*pbMsg.BuyerId, pbMsg.PurchasedGigs[0], pbMsg.Action.String())
				} else {
					slog.Debug("Buyer Consumer: missing purchasedGig information")
				}
			}

			msg.Ack(false)
		}

		slog.Debug("ConsumeBuyerDirectMessage shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}

func (q *Queue) ConsumeNotificationRelay() error {
	var (
		exchangeName = "jobber-relay-notification"
		routingKey   = "relay-notification"
		queueName    = "user-relay-queue"
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
			var pbMsg pbEmail.OrderEmailMessageDetails

			if err := proto.Unmarshal(msg.Body, &pbMsg); err != nil {
				slog.With("error", err).Error("ConsumeNotificationRelay: failed to unmarshal protobuf message")
				continue
			}

			// get the buyer by username and update the receiver email with the buyer email address.
			buyer, err := q.buyer.GetBuyerByUsername(context.Background(), *pbMsg.BuyerUsername)
			if err != nil {
				slog.With("error", err).Error("failed to get buyer email")

				msg.Ack(false)
				continue
			}

			pbMsg.ReceiverEmail = &buyer.Email

			// send the offer information to the notification micro service via rabbitMQ direct exchange using protobuf
			if data, err := proto.Marshal(&pbMsg); err == nil {
				if err := q.PublishDirectMessage("jobber-order-notification", "order-email", data); err != nil {
					slog.With("error", err).Error("ConsumeNotificationRelay: Failed to send message to jobber-order-notification")
				}
			}

			msg.Ack(false)
		}

		slog.Debug("ConsumeNotificationRelay shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}
