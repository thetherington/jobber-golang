package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/order"
	"github.com/thetherington/jobber-common/models/review"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *order.OrderDocument) (*order.OrderResponse, error)
	CreatePaymentIntent(ctx context.Context, price float32, buyerId string) (*order.PaymentIntentResponse, error)
	GetOrderById(ctx context.Context, id string) (*order.OrderResponse, error)
	GetSellerOrders(ctx context.Context, id string) (*order.OrdersResponse, error)
	GetBuyerOrders(ctx context.Context, id string) (*order.OrdersResponse, error)
	CancelOrder(ctx context.Context, id string, paymentIntentId string, orderData *order.OrderMessage) (*order.OrderResponse, error)
	RequestExtension(ctx context.Context, id string, extension *order.ExtendedDelivery) (*order.OrderResponse, error)
	DeliverOrder(ctx context.Context, id string, file *order.DeliveredWork) (*order.OrderResponse, error)
	ApproveOrder(ctx context.Context, id string, orderData *order.OrderMessage) (*order.OrderResponse, error)
	RejectDeliveryDate(ctx context.Context, orderId string) (*order.OrderResponse, error)
	ApproveDeliveryDate(ctx context.Context, orderId string, data order.ExtendedDelivery) (*order.OrderResponse, error)
	UpdateOrderReview(ctx context.Context, data *review.ReviewMessageDetails) (*order.OrderResponse, string, error)
}

type GrpcInterface interface {
	NotifyUpdateOrder(order *order.OrderDocument, userTo string, message string)
}

type ImageUploader interface {
	UploadImage(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
	UploadVideo(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error)
}

type OrderProducer interface {
	PublishDirectMessage(exchangeName string, routingKey string, data []byte) error
}

type NotificationService interface {
	CreateNotification(data *order.Notification) (*order.Notification, error)
	GetNotificationsById(ctx context.Context, userToId string) ([]*order.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationId string) (*order.Notification, error)
	CreateReturnNotification(order *order.OrderDocument, userToId string, message string) (*order.Notification, error)
}

type PaymentService interface {
	SearchCustomers(email string) (string, error)
	CreateCustomer(email string, buyerId string) (string, error)
	CreatePaymentIntent(customerId string, price float32) (string, string, error)
	RefundOrder(paymentIntent string) error
}
