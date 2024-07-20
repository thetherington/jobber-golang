package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/order"
	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"google.golang.org/grpc"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *order.OrderDocument) (*order.OrderResponse, error)
	CreatePaymentIntent(ctx context.Context, price float32, buyerId string) (*order.PaymentIntentResponse, error)
	GetOrderById(ctx context.Context, id string) (*order.OrderResponse, error)
	GetSellerOrders(ctx context.Context, id string) (*order.OrdersResponse, error)
	GetBuyerOrders(ctx context.Context, id string) (*order.OrdersResponse, error)
	CancelOrder(ctx context.Context, id string, paymentIntentId string, orderData *order.OrderMessage) (string, error)
	RequestExtension(ctx context.Context, id string, extension *order.ExtendedDelivery) (*order.OrderResponse, error)
	DeliverOrder(ctx context.Context, id string, file *order.DeliveredWork) (*order.OrderResponse, error)
	ApproveOrder(ctx context.Context, id string, orderData *order.OrderMessage) (*order.OrderResponse, error)
	DeliveryDate(ctx context.Context, id string, deliveryChangeAction string, extenstion *order.ExtendedDelivery) (*order.OrderResponse, error)
	GetNotificationsById(ctx context.Context, id string) (*order.NotificationsResponse, error)
	MarkNotificationAsRead(ctx context.Context, id string) (*order.NotificationResponse, error)
}

type OrderRPCClient interface {
	Create(ctx context.Context, in *pb.OrderDocument, opts ...grpc.CallOption) (*pb.OrderResponse, error)
	Intent(ctx context.Context, in *pb.PaymentIntentRequest, opts ...grpc.CallOption) (*pb.PaymentIntentResponse, error)
	GetOrderById(ctx context.Context, in *pb.RequestById, opts ...grpc.CallOption) (*pb.OrderResponse, error)
	GetSellerOrders(ctx context.Context, in *pb.RequestById, opts ...grpc.CallOption) (*pb.OrdersResponse, error)
	GetBuyerOrders(ctx context.Context, in *pb.RequestById, opts ...grpc.CallOption) (*pb.OrdersResponse, error)
	RequestExtension(ctx context.Context, in *pb.DeliveryExtensionRequest, opts ...grpc.CallOption) (*pb.OrderResponse, error)
	DeliveryDate(ctx context.Context, in *pb.DeliveryDateRequest, opts ...grpc.CallOption) (*pb.OrderResponse, error)
	DeliverOrder(ctx context.Context, in *pb.DeliverOrderRequest, opts ...grpc.CallOption) (*pb.OrderResponse, error)
	CancelOrder(ctx context.Context, in *pb.CancelOrderRequest, opts ...grpc.CallOption) (*pb.MessageResponse, error)
	ApproveOrder(ctx context.Context, in *pb.ApproveOrderRequest, opts ...grpc.CallOption) (*pb.OrderResponse, error)
	GetNotificationsById(ctx context.Context, in *pb.RequestWithParam, opts ...grpc.CallOption) (*pb.NotificationsResponse, error)
	MarkNotificationAsRead(ctx context.Context, in *pb.RequestWithParam, opts ...grpc.CallOption) (*pb.NotificationResponse, error)
}

type NotificationSocketManager interface {
	DispatchNotification(cmd string, payload *order.Notification)
	DispatchOrder(cmd string, payload *order.OrderDocument)
}
