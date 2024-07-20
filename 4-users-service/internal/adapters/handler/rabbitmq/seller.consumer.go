package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thetherington/jobber-common/models/review"
	"github.com/thetherington/jobber-common/models/users"
	pborder "github.com/thetherington/jobber-common/protogen/go/order"
	pbReview "github.com/thetherington/jobber-common/protogen/go/review"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"google.golang.org/protobuf/proto"
)

func (q *Queue) ConsumeSellerDirectMessage() error {
	var (
		exchangeName = "jobber-seller-update"
		routingKey   = "user-seller"
		queueName    = "user-seller-queue"
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
			var order pborder.SellerGigUpdate

			if err := proto.Unmarshal(msg.Body, &order); err != nil {
				slog.With("error", err).Error("ConsumeSellerDirectMessage: failed to unmarshal protobuf message")
				continue
			}

			switch order.Action {
			case pborder.Action_CreateOrder:
				q.seller.UpdateSellerOngoingJobsProp(order.SellerId, *order.OrderProps.OngoingJobs)

			case pborder.Action_ApproveOrder:
				q.seller.UpdateSellerCompletedJobsProp(
					order.SellerId,
					*order.OrderProps.OngoingJobs,
					*order.OrderProps.CompletedJobs,
					*order.OrderProps.TotalEarnings,
				)

			case pborder.Action_UpdateGigCount:
				q.seller.UpdateTotalGigCount(order.SellerId, *order.OrderProps.GigCount)

			case pborder.Action_CancelOrder:
				q.seller.UpdateSellerCancelledJobsProp(order.SellerId)

			default:
				slog.Info("Unknown event type", "type", order.Action.String())
			}

			msg.Ack(false)
		}

		slog.Debug("ConsumeSellerDirectMessage shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}

// Fanout exchange
func (q *Queue) ConsumeReviewFanoutMessages() error {
	var (
		exchangeName = "jobber-review"
		queueName    = "seller-review-queue"
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
				slog.With("error", err).Error("ConsumeReviewFanoutMessages: failed to unmarshal protobuf message")
				continue
			}

			switch pbMsg.Action {
			case pbReview.ReviewType_BuyerReview:
				err := q.seller.UpdateSellerReview(&review.ReviewMessageDetails{
					GigId:      pbMsg.GigId,
					ReviewerId: pbMsg.ReviewerId,
					SellerId:   pbMsg.SellerId,
					Review:     pbMsg.Review,
					Rating:     pbMsg.Rating,
					OrderId:    pbMsg.OrderId,
					Type:       pbMsg.Action.String(),
				})
				if err != nil {
					slog.With("error", err).Error("failed to update seller review")
				}

				pbMsg.Action = pbReview.ReviewType_GigReview

				// send message to the gig service via the gig exchange
				if data, err := proto.Marshal(&pbMsg); err == nil {
					if err := q.PublishDirectMessage("jobber-update-gig", "update-gig", data); err != nil {
						slog.With("error", err).Error("Review: Failed to send review data to jobber-update-gig")
					}
				}
			}

			msg.Ack(false)
		}

		slog.Debug("ConsumeSellerDirectMessage shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}

func (q *Queue) ConsumeSeedGigDirectMessages() error {
	var (
		exchangeName = "jobber-gig"
		routingKey   = "get-sellers"
		queueName    = "user-gig-queue"
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
			var pbMsg pb.RandomSellersRequest

			if err := proto.Unmarshal(msg.Body, &pbMsg); err != nil {
				slog.With("error", err).Error("ConsumeSeedGigDirectMessages: failed to unmarshal protobuf message")
				continue
			}

			if resp, err := q.seller.GetRandomSellers(context.Background(), pbMsg.Size); err == nil {
				sellers := make([]*pb.SellerPayload, 0)
				for _, s := range resp.Sellers {
					sellers = append(sellers, users.CreatePayloadFromSeller(s))
				}

				pbSellers := &pb.SellersResponse{
					Message: "random sellers",
					Sellers: sellers,
				}

				// send message to the gig service via the gig exchange
				if data, err := proto.Marshal(pbSellers); err == nil {
					if err := q.PublishDirectMessage("jobber-seed-gig", "receive-sellers", data); err != nil {
						slog.With("error", err).Error("Seed: Failed to send sellers to jobber-seed-gig")
					}
				}
			}

			msg.Ack(false)
		}

		slog.Debug("ConsumeSeedGigDirectMessages shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}
