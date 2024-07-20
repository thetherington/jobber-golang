package grpc

import (
	"context"

	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-common/utils"
)

func (g *GrpcAdapter) GetNotificationsById(ctx context.Context, req *pb.RequestWithParam) (*pb.NotificationsResponse, error) {
	resp, err := g.notificationService.GetNotificationsById(ctx, req.Param)
	if err != nil {
		return nil, serviceError(err)
	}

	notifications := make([]*pb.NotificationMessage, 0)

	for _, n := range resp {
		notifications = append(notifications, &pb.NotificationMessage{
			Id:               n.Id,
			UserTo:           n.UserTo,
			SenderUsername:   n.SenderUsername,
			SenderPicture:    n.SenderPicture,
			ReceiverUsername: n.ReceiverUsername,
			ReceiverPicture:  n.ReceiverPicture,
			IsRead:           n.IsRead,
			Message:          n.Message,
			OrderId:          n.OrderId,
			Cmd:              "",
			CreatedAt:        utils.ToDateTime(n.CreatedAt),
		})
	}

	return &pb.NotificationsResponse{Message: "Notifications", Notifications: notifications}, nil
}

func (g *GrpcAdapter) MarkNotificationAsRead(ctx context.Context, req *pb.RequestWithParam) (*pb.NotificationResponse, error) {
	resp, err := g.notificationService.MarkNotificationAsRead(ctx, req.Param)
	if err != nil {
		return nil, serviceError(err)
	}

	notification := &pb.NotificationMessage{
		Id:               resp.Id,
		UserTo:           resp.UserTo,
		SenderUsername:   resp.SenderUsername,
		SenderPicture:    resp.SenderPicture,
		ReceiverUsername: resp.ReceiverUsername,
		ReceiverPicture:  resp.ReceiverPicture,
		IsRead:           resp.IsRead,
		Message:          resp.Message,
		OrderId:          resp.OrderId,
		Cmd:              "",
		CreatedAt:        utils.ToDateTime(resp.CreatedAt),
	}

	return &pb.NotificationResponse{Message: "Notification updated successfully.", Notification: notification}, nil
}
