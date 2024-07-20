package service

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/order"
	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

/**
 * OrderService implements
 */
type OrderService struct {
	client port.OrderRPCClient
}

// NewOrderService creates a new order service instance
func NewOrderService(rpc port.OrderRPCClient) *OrderService {
	return &OrderService{
		rpc,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, payload *order.OrderDocument) (*order.OrderResponse, error) {
	resp, err := o.client.Create(ctx, payload.MarshalToProto())
	if err != nil {
		slog.With("error", err).Debug("CreateGig error")
		return nil, svc.GrpcErrorResolve(err, "CreateOrder")
	}

	return &order.OrderResponse{Message: resp.Message, Order: order.CreateOrderDocument(resp.Order)}, nil
}

func (o *OrderService) CreatePaymentIntent(ctx context.Context, price float32, buyerId string) (*order.PaymentIntentResponse, error) {
	resp, err := o.client.Intent(ctx, &pb.PaymentIntentRequest{Price: price, BuyerId: buyerId})
	if err != nil {
		slog.With("error", err).Debug("CreatePaymentIntent error")
		return nil, svc.GrpcErrorResolve(err, "CreatePaymentIntent")
	}

	return &order.PaymentIntentResponse{
		Message: resp.Message, ClientSecret: resp.ClientSecret, PaymentIntentId: resp.PaymentIntentId,
	}, nil
}

func (o *OrderService) GetOrderById(ctx context.Context, id string) (*order.OrderResponse, error) {
	resp, err := o.client.GetOrderById(ctx, &pb.RequestById{Id: id})
	if err != nil {
		slog.With("error", err).Debug("GetOrderById error")
		return nil, svc.GrpcErrorResolve(err, "GetOrderById")
	}

	return &order.OrderResponse{Message: resp.Message, Order: order.CreateOrderDocument(resp.Order)}, nil
}

func (o *OrderService) GetSellerOrders(ctx context.Context, id string) (*order.OrdersResponse, error) {
	resp, err := o.client.GetSellerOrders(ctx, &pb.RequestById{Id: id})
	if err != nil {
		slog.With("error", err).Debug("GetSellerOrders error")
		return nil, svc.GrpcErrorResolve(err, "GetSellerOrders")
	}

	orders := make([]*order.OrderDocument, 0)

	for _, o := range resp.Orders {
		orders = append(orders, order.CreateOrderDocument(o))
	}

	return &order.OrdersResponse{Message: resp.Message, Orders: orders}, nil
}

func (o *OrderService) GetBuyerOrders(ctx context.Context, id string) (*order.OrdersResponse, error) {
	resp, err := o.client.GetBuyerOrders(ctx, &pb.RequestById{Id: id})
	if err != nil {
		slog.With("error", err).Debug("GetBuyerOrders error")
		return nil, svc.GrpcErrorResolve(err, "GetBuyerOrders")
	}

	orders := make([]*order.OrderDocument, 0)

	for _, o := range resp.Orders {
		orders = append(orders, order.CreateOrderDocument(o))
	}

	return &order.OrdersResponse{Message: resp.Message, Orders: orders}, nil
}

func (o *OrderService) CancelOrder(ctx context.Context, id string, paymentIntentId string, orderData *order.OrderMessage) (string, error) {
	resp, err := o.client.CancelOrder(ctx, &pb.CancelOrderRequest{
		OrderId: id, PaymentIntent: paymentIntentId, OrderMessage: orderData.MarshalToProto(),
	})
	if err != nil {
		slog.With("error", err).Debug("CancelOrder error")
		return "", svc.GrpcErrorResolve(err, "CancelOrder")
	}

	return resp.Message, nil
}

func (o *OrderService) RequestExtension(ctx context.Context, id string, extension *order.ExtendedDelivery) (*order.OrderResponse, error) {
	req := &pb.DeliveryExtensionRequest{
		OrderId: id,
		ExtendedDelivery: &pb.ExtendedDelivery{
			OriginalDate:        utils.ToDateTimeOrNil(extension.OriginalDate),
			NewDate:             utils.ToDateTimeOrNil(extension.NewDate),
			Days:                extension.Days,
			Reason:              extension.Reason,
			DeliveryDateUpdated: extension.DeliveryDateUpdate,
		},
	}

	resp, err := o.client.RequestExtension(ctx, req)
	if err != nil {
		slog.With("error", err).Debug("RequestExtension error")
		return nil, svc.GrpcErrorResolve(err, "RequestExtension")
	}

	return &order.OrderResponse{Message: resp.Message, Order: order.CreateOrderDocument(resp.Order)}, nil
}

func (o *OrderService) DeliverOrder(ctx context.Context, id string, file *order.DeliveredWork) (*order.OrderResponse, error) {
	req := &pb.DeliverOrderRequest{
		OrderId:   id,
		Delivered: true,
		DeliveredWork: &pb.DeliveredWork{
			Message:  file.Message,
			File:     file.File,
			FileType: file.FileType,
			FileSize: file.FileSize,
			FileName: file.FileName,
		},
	}

	resp, err := o.client.DeliverOrder(ctx, req)
	if err != nil {
		slog.With("error", err).Debug("DeliverOrder error")
		return nil, svc.GrpcErrorResolve(err, "DeliverOrder")
	}

	return &order.OrderResponse{Message: resp.Message, Order: order.CreateOrderDocument(resp.Order)}, nil
}

func (o *OrderService) ApproveOrder(ctx context.Context, id string, orderData *order.OrderMessage) (*order.OrderResponse, error) {
	resp, err := o.client.ApproveOrder(ctx, &pb.ApproveOrderRequest{
		OrderId:      id,
		OrderMessage: orderData.MarshalToProto(),
	})
	if err != nil {
		slog.With("error", err).Debug("DeliverOrder error")
		return nil, svc.GrpcErrorResolve(err, "DeliverOrder")
	}

	return &order.OrderResponse{Message: resp.Message, Order: order.CreateOrderDocument(resp.Order)}, nil
}

func (o *OrderService) DeliveryDate(ctx context.Context, id string, deliveryChangeAction string, extenstion *order.ExtendedDelivery) (*order.OrderResponse, error) {
	var action pb.DeliveryChangeAction

	switch deliveryChangeAction {
	case "approve":
		action = pb.DeliveryChangeAction_Approve
	case "reject":
		action = pb.DeliveryChangeAction_Reject
	}

	req := &pb.DeliveryDateRequest{
		OrderId: id,
		Action:  action,
		RequestExtension: &pb.ExtendedDelivery{
			OriginalDate:        utils.ToDateTime(extenstion.OriginalDate),
			NewDate:             utils.ToDateTime(extenstion.NewDate),
			Days:                extenstion.Days,
			Reason:              extenstion.Reason,
			DeliveryDateUpdated: extenstion.DeliveryDateUpdate,
		},
	}

	resp, err := o.client.DeliveryDate(ctx, req)
	if err != nil {
		slog.With("error", err).Debug("DeliveryDate error")
		return nil, svc.GrpcErrorResolve(err, "DeliveryDate")
	}

	return &order.OrderResponse{Message: resp.Message, Order: order.CreateOrderDocument(resp.Order)}, nil
}

func (o *OrderService) GetNotificationsById(ctx context.Context, id string) (*order.NotificationsResponse, error) {
	resp, err := o.client.GetNotificationsById(ctx, &pb.RequestWithParam{Param: id})
	if err != nil {
		slog.With("error", err).Debug("GetNotificationsById error")
		return nil, svc.GrpcErrorResolve(err, "GetNotificationsById")
	}

	notifcations := make([]*order.Notification, 0)

	for _, n := range resp.Notifications {
		notifcations = append(notifcations, &order.Notification{
			Id:               n.Id,
			UserTo:           n.UserTo,
			SenderUsername:   n.SenderUsername,
			SenderPicture:    n.SenderPicture,
			ReceiverUsername: n.ReceiverUsername,
			ReceiverPicture:  n.ReceiverPicture,
			IsRead:           n.IsRead,
			Message:          n.Message,
			OrderId:          n.OrderId,
			CreatedAt:        utils.ToTime(n.CreatedAt),
		})
	}

	return &order.NotificationsResponse{Message: resp.Message, Notifications: notifcations}, nil
}

func (o *OrderService) MarkNotificationAsRead(ctx context.Context, id string) (*order.NotificationResponse, error) {
	resp, err := o.client.MarkNotificationAsRead(ctx, &pb.RequestWithParam{Param: id})
	if err != nil {
		slog.With("error", err).Debug("MarkNotificationAsRead error")
		return nil, svc.GrpcErrorResolve(err, "MarkNotificationAsRead")
	}

	n := &order.Notification{
		Id:               resp.Notification.Id,
		UserTo:           resp.Notification.UserTo,
		SenderUsername:   resp.Notification.SenderUsername,
		SenderPicture:    resp.Notification.SenderPicture,
		ReceiverUsername: resp.Notification.ReceiverUsername,
		ReceiverPicture:  resp.Notification.ReceiverPicture,
		IsRead:           resp.Notification.IsRead,
		Message:          resp.Notification.Message,
		OrderId:          resp.Notification.OrderId,
		CreatedAt:        utils.ToTime(resp.Notification.CreatedAt),
	}

	return &order.NotificationResponse{Message: resp.Message, Notification: n}, nil
}
