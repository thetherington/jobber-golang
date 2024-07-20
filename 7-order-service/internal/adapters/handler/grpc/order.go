package grpc

import (
	"context"
	"log/slog"

	"github.com/thetherington/jobber-common/error-handling/grpcerror"
	"github.com/thetherington/jobber-common/models/order"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/thetherington/jobber-common/protogen/go/order"
)

// parses service error and returns a equivelient gRPC error
func serviceError(err error) error {
	// try to cast the error to a grpcerror lookup
	if apiError, ok := grpcerror.FromError(err); ok {
		s := status.New(apiError.Status, apiError.Message)
		return s.Err()
	}

	// generic response
	s := status.New(codes.Internal, err.Error())
	return s.Err()
}

func (g *GrpcAdapter) Create(ctx context.Context, req *pb.OrderDocument) (*pb.OrderResponse, error) {
	resp, err := g.orderService.CreateOrder(ctx, order.CreateOrderDocument(req))
	if err != nil {
		return nil, serviceError(err)
	}

	// send notification
	go g.NotifyUpdateOrder(resp.Order, resp.Order.SellerUsername, "placed an order for your gig.")

	return &pb.OrderResponse{Message: resp.Message, Order: resp.Order.MarshalToProto()}, nil
}

func (g *GrpcAdapter) Intent(ctx context.Context, req *pb.PaymentIntentRequest) (*pb.PaymentIntentResponse, error) {
	resp, err := g.orderService.CreatePaymentIntent(ctx, req.Price, req.BuyerId)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.PaymentIntentResponse{
		Message:         resp.Message,
		ClientSecret:    resp.ClientSecret,
		PaymentIntentId: resp.PaymentIntentId,
	}, nil
}

func (g *GrpcAdapter) GetOrderById(ctx context.Context, req *pb.RequestById) (*pb.OrderResponse, error) {
	resp, err := g.orderService.GetOrderById(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	return &pb.OrderResponse{Message: resp.Message, Order: resp.Order.MarshalToProto()}, nil
}

func (g *GrpcAdapter) GetSellerOrders(ctx context.Context, req *pb.RequestById) (*pb.OrdersResponse, error) {
	resp, err := g.orderService.GetSellerOrders(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	orders := make([]*pb.OrderDocument, 0)

	for _, o := range resp.Orders {
		orders = append(orders, o.MarshalToProto())
	}

	return &pb.OrdersResponse{Message: resp.Message, Orders: orders}, nil
}

func (g *GrpcAdapter) GetBuyerOrders(ctx context.Context, req *pb.RequestById) (*pb.OrdersResponse, error) {
	resp, err := g.orderService.GetBuyerOrders(ctx, req.Id)
	if err != nil {
		return nil, serviceError(err)
	}

	orders := make([]*pb.OrderDocument, 0)

	for _, o := range resp.Orders {
		orders = append(orders, o.MarshalToProto())
	}

	return &pb.OrdersResponse{Message: resp.Message, Orders: orders}, nil
}

func (g *GrpcAdapter) RequestExtension(ctx context.Context, req *pb.DeliveryExtensionRequest) (*pb.OrderResponse, error) {
	resp, err := g.orderService.RequestExtension(ctx, req.OrderId, &order.ExtendedDelivery{
		OriginalDate:       utils.ToTimeOrNil(req.ExtendedDelivery.OriginalDate),
		NewDate:            utils.ToTimeOrNil(req.ExtendedDelivery.NewDate),
		Days:               req.ExtendedDelivery.Days,
		Reason:             req.ExtendedDelivery.Reason,
		DeliveryDateUpdate: req.ExtendedDelivery.DeliveryDateUpdated,
	})
	if err != nil {
		return nil, serviceError(err)
	}

	// send notification
	go g.NotifyUpdateOrder(resp.Order, resp.Order.BuyerUsername, "requested for an order delivery date extension")

	return &pb.OrderResponse{Message: resp.Message, Order: resp.Order.MarshalToProto()}, nil
}

func (g *GrpcAdapter) DeliveryDate(ctx context.Context, req *pb.DeliveryDateRequest) (*pb.OrderResponse, error) {
	var (
		resp *order.OrderResponse
		err  error
	)

	// approved order extenstion
	if req.Action == pb.DeliveryChangeAction_Approve {
		resp, err = g.orderService.ApproveDeliveryDate(ctx, req.OrderId, order.ExtendedDelivery{
			OriginalDate:       utils.ToTimeOrNil(req.RequestExtension.OriginalDate),
			NewDate:            utils.ToTimeOrNil(req.RequestExtension.NewDate),
			Days:               req.RequestExtension.Days,
			Reason:             req.RequestExtension.Reason,
			DeliveryDateUpdate: req.RequestExtension.DeliveryDateUpdated,
		})
		if err != nil {
			return nil, serviceError(err)
		}

		// send notification
		go g.NotifyUpdateOrder(resp.Order, resp.Order.SellerUsername, "approved your order delivery date extension request.")
	}

	// rejected order extenstion
	if req.Action == pb.DeliveryChangeAction_Reject {
		resp, err = g.orderService.RejectDeliveryDate(ctx, req.OrderId)
		if err != nil {
			return nil, serviceError(err)
		}

		// send notification
		go g.NotifyUpdateOrder(resp.Order, resp.Order.SellerUsername, "rejected your order delivery date extension request.")
	}

	return &pb.OrderResponse{Message: resp.Message, Order: resp.Order.MarshalToProto()}, nil
}

func (g *GrpcAdapter) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.MessageResponse, error) {
	resp, err := g.orderService.CancelOrder(ctx, req.OrderId, req.PaymentIntent, order.CreateOrderMessage(req.OrderMessage))
	if err != nil {
		return nil, serviceError(err)
	}

	// send notification
	go g.NotifyUpdateOrder(resp.Order, resp.Order.BuyerUsername, "order was cancelled.")

	return &pb.MessageResponse{Message: resp.Message}, nil
}

func (g *GrpcAdapter) ApproveOrder(ctx context.Context, req *pb.ApproveOrderRequest) (*pb.OrderResponse, error) {
	resp, err := g.orderService.ApproveOrder(ctx, req.OrderId, order.CreateOrderMessage(req.OrderMessage))
	if err != nil {
		return nil, serviceError(err)
	}

	// send notification
	go g.NotifyUpdateOrder(resp.Order, resp.Order.SellerUsername, "approved your order delivery.")

	return &pb.OrderResponse{Message: resp.Message, Order: resp.Order.MarshalToProto()}, nil
}

func (g *GrpcAdapter) DeliverOrder(ctx context.Context, req *pb.DeliverOrderRequest) (*pb.OrderResponse, error) {
	resp, err := g.orderService.DeliverOrder(ctx, req.OrderId, &order.DeliveredWork{
		Message:  req.DeliveredWork.Message,
		File:     req.DeliveredWork.File,
		FileType: req.DeliveredWork.FileType,
		FileSize: req.DeliveredWork.FileSize,
		FileName: req.DeliveredWork.FileName,
	})
	if err != nil {
		return nil, serviceError(err)
	}

	// send notification
	go g.NotifyUpdateOrder(resp.Order, resp.Order.BuyerUsername, "delivered your order.")

	return &pb.OrderResponse{Message: resp.Message, Order: resp.Order.MarshalToProto()}, nil
}

func (g *GrpcAdapter) NotifyUpdateOrder(order *order.OrderDocument, userTo string, message string) {
	// create notification message
	n, err := g.notificationService.CreateReturnNotification(order, userTo, message)
	if err != nil {
		slog.With("err", err).Error("failed to generate notification for order", "id", order.OrderId)
		return
	}

	// publish notification on return streaming
	g.PushMessage("order_notification", n)

	// push order update
	g.PushOrder("order_update", order)
}
