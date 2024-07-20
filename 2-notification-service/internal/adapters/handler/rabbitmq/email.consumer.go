package rabbitmq

import (
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thetherington/jobber-common/models/notification"
	pb "github.com/thetherington/jobber-common/protogen/go/notification"
	"google.golang.org/protobuf/proto"
)

func (q *Queue) ConsumeAuthEmailMessages() error {
	var (
		exchangeName = "jobber-email-notification"
		routingKey   = "auth-email"
		queueName    = "auth-email-queue"
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
			var pbMsg pb.EmailMessageDetails

			if err := proto.Unmarshal(msg.Body, &pbMsg); err != nil {
				slog.With("error", err).Error("Email Consumer: failed to unmarshal protobuf message")
				continue
			}

			locals := notification.AuthEmailLocals{
				Username:      pbMsg.GetUsername(),
				VerifyLink:    pbMsg.GetVerifyLink(),
				ResetLink:     pbMsg.GetResetLink(),
				ReceiverEmail: pbMsg.GetReceiverEmail(),
			}

			// send email via service using a go routine
			q.service.SendAuthEmail(pbMsg.GetTemplate(), locals)

			msg.Ack(false)
		}

		slog.Debug("ConsumeAuthEmailMessages shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}

func (q *Queue) ConsumeOrderEmailMessages() error {
	var (
		exchangeName = "jobber-order-notification"
		routingKey   = "order-email"
		queueName    = "order-email-queue"
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
			var pbMsg pb.OrderEmailMessageDetails

			if err := proto.Unmarshal(msg.Body, &pbMsg); err != nil {
				slog.With("error", err).Error("ConsumeOrderEmailMessages: failed to unmarshal protobuf message")
				continue
			}

			locals := notification.OrderEmailLocals{
				ReceiverEmail:  pbMsg.GetReceiverEmail(),
				Username:       pbMsg.GetUsername(),
				Template:       pbMsg.GetTemplate(),
				Sender:         pbMsg.GetSender(),
				OfferLink:      pbMsg.GetOfferLink(),
				Amount:         pbMsg.GetAmount(),
				BuyerUsername:  pbMsg.GetBuyerUsername(),
				SellerUsername: pbMsg.GetSellerUsername(),
				Title:          pbMsg.GetTitle(),
				Description:    pbMsg.GetDescription(),
				DeliveryDays:   pbMsg.GetDeliveryDays(),
				InvoiceId:      pbMsg.GetInvoiceId(),
				OrderId:        pbMsg.GetOrderId(),
				OrderDue:       pbMsg.GetOrderDue(),
				Requirements:   pbMsg.GetRequirements(),
				OrderUrl:       pbMsg.GetOrderUrl(),
				OriginalDate:   pbMsg.GetOriginalDate(),
				NewDate:        pbMsg.GetNewDate(),
				Reason:         pbMsg.GetReason(),
				Subject:        pbMsg.GetSubject(),
				Header:         pbMsg.GetHeader(),
				Type:           pbMsg.GetType(),
				Message:        pbMsg.GetMessage(),
				ServiceFee:     pbMsg.GetServiceFee(),
				Total:          pbMsg.GetTotal(),
			}

			// send email via service using a go routine? using a go routine causes emails to fail for some reason
			q.service.SendOrderEmail(pbMsg.GetTemplate(), locals)

			msg.Ack(false)
		}

		slog.Debug("ConsumeOrderEmailMessages shutting down")
		ch.Close()
	}()

	slog.Info("Waiting for messages", "[Exchange,Queue]", fmt.Sprintf("[%s,%s]", exchangeName, queueName))

	return nil
}
