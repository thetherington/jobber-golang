package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/middleware"
	"github.com/thetherington/jobber-common/models/order"
	pbnotify "github.com/thetherington/jobber-common/protogen/go/notification"
	pborder "github.com/thetherington/jobber-common/protogen/go/order"
	pbBuyer "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-order/internal/adapters/config"
	"github.com/thetherington/jobber-order/internal/core/port"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

var (
	orderDB  *mongo.Collection
	validate *validator.Validate
)

/**
 * OrderService implements
 */
type OrderService struct {
	queue   port.OrderProducer
	image   port.ImageUploader
	payment port.PaymentService
}

// NewOrderService creates a new order service instance
func NewOrderService(db *mongo.Database, queue port.OrderProducer, image port.ImageUploader, payment port.PaymentService) *OrderService {
	validate = validator.New(validator.WithRequiredStructEnabled())

	orderDB = db.Collection("Order")

	orderDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "orderId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "buyerId", Value: 1}}},
		{Keys: bson.D{{Key: "sellerId", Value: 1}}},
	})

	return &OrderService{
		queue:   queue,
		image:   image,
		payment: payment,
	}
}

func (o *OrderService) SetQueue(queue port.OrderProducer) {
	o.queue = queue
}

func (o *OrderService) CreateOrder(ctx context.Context, payload *order.OrderDocument) (*order.OrderResponse, error) {
	// Validate Order payload
	if err := payload.Validate(validate); err != nil {
		slog.With("error", err).Debug("Order Create Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// Validate offer in order payload
	if err := payload.Offer.Validate(validate); err != nil {
		slog.With("error", err).Debug("Order Create Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// the service charge is 5.5% of the purchase amount
	// for purchases under $50, an additional $2 is applied
	var serviceFee float64

	if payload.Price < 50 {
		serviceFee, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", ((5.5/100)*payload.Price)+2), 32)
	} else {
		serviceFee, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (5.5/100)*payload.Price), 32)
	}

	payload.ServiceFee = float32(serviceFee)

	// save order into the database
	_, err := orderDB.InsertOne(ctx, payload)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// publish message for update seller from users service
	pbSellerUpdate := &pborder.SellerGigUpdate{
		Action:   pborder.Action_CreateOrder,
		SellerId: payload.SellerId,
		OrderProps: &pborder.OrderProps{
			OngoingJobs: proto.Int32(1),
		},
	}
	// send 1 ongoing job count to the users microservice to increment the seller ongoing jobs count
	if data, err := proto.Marshal(pbSellerUpdate); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-seller-update", "user-seller", data); err != nil {
			slog.With("error", err).Error("CreateOrder: Failed to send offer to jobber-seller-update")
		}
	}

	// publish message send seller email
	pbEmail := &pbnotify.OrderEmailMessageDetails{
		ReceiverEmail:  &payload.SellerEmail,
		Template:       utils.Ptr("orderPlaced"),
		Amount:         utils.Ptr(fmt.Sprintf("%.2f", payload.Price)),
		BuyerUsername:  &payload.BuyerUsername,
		SellerUsername: &payload.SellerUsername,
		Title:          &payload.Offer.GigTitle,
		Description:    &payload.Offer.Description,
		InvoiceId:      &payload.InvoiceId,
		OrderId:        &payload.OrderId,
		OrderDue:       &payload.Offer.NewDeliveryDate,
		Requirements:   &payload.Requirements,
		OrderUrl:       utils.Ptr(fmt.Sprintf("%s/orders/%s/activities", config.Config.App.ClientUrl, payload.OrderId)),
		ServiceFee:     utils.Ptr(fmt.Sprintf("%.2f", payload.ServiceFee)),
		Total:          utils.Ptr(fmt.Sprintf("%.2f", payload.ServiceFee+payload.Price)),
	}

	// send the order information to the notification microservice to email the seller
	if data, err := proto.Marshal(pbEmail); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-order-notification", "order-email", data); err != nil {
			slog.With("error", err).Error("CreateOrder: Failed to send offer to jobber-order-notification")
		}
	}

	// publish message send buyer email
	pbEmail.ReceiverEmail = &payload.BuyerEmail
	pbEmail.Template = utils.Ptr("orderReceipt")

	// send the order information to the notification microservice to email the buyer
	if data, err := proto.Marshal(pbEmail); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-order-notification", "order-email", data); err != nil {
			slog.With("error", err).Error("CreateOrder: Failed to send offer to jobber-order-notification")
		}
	}

	// update buyer orders from User microservice
	pbBuyer := &pbBuyer.BuyerPayload{
		Action:        pbBuyer.Action_PURCHASED_GIG.Enum(),
		BuyerId:       utils.Ptr(payload.BuyerId),
		PurchasedGigs: []string{payload.GigId},
	}

	if data, err := proto.Marshal(pbBuyer); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-buyer-update", "user-buyer", data); err != nil {
			slog.With("error", err).Error("CreateOrder: Failed to send offer to jobber-buyer-update")
		}
	}

	return &order.OrderResponse{Message: "Order created successfully.", Order: payload}, nil
}

func (o *OrderService) CreatePaymentIntent(ctx context.Context, price float32, buyerId string) (*order.PaymentIntentResponse, error) {
	// Get email from the user cookie session passed down into the context.
	email := ctx.Value(middleware.CtxEmailKey)
	if email == nil {
		slog.Debug("Email in context is nil")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again"))
	}

	var (
		customerId string
		err        error
	)

	// search for customers with email address
	customerId, err = o.payment.SearchCustomers(email.(string))
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// create stripe customer if not found
	if customerId == "" {
		customerId, err = o.payment.CreateCustomer(email.(string), buyerId)
		if err != nil {
			return nil, svc.NewError(svc.ErrInternalFailure, err)
		}
	}

	// if customerId is still blank, then return error
	if customerId == "" {
		return nil, svc.NewError(svc.ErrInternalFailure, fmt.Errorf("no valid stripe customer information"))
	}

	// the service charge is 5.5% of the purchase amount
	// for purchases under $50, an additional $2 is applied
	var serviceFee float64

	if price < 50 {
		serviceFee, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", ((5.5/100)*price)+2), 32)
	} else {
		serviceFee, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (5.5/100)*price), 32)
	}

	// create the payment intent
	paymentIntentId, clientSecret, err := o.payment.CreatePaymentIntent(customerId, price+float32(serviceFee))
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &order.PaymentIntentResponse{
		Message:         "Order intent created successfully",
		ClientSecret:    clientSecret,
		PaymentIntentId: paymentIntentId,
	}, nil
}
