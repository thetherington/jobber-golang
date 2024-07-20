package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/order"
	"github.com/thetherington/jobber-common/models/review"
	pbnotify "github.com/thetherington/jobber-common/protogen/go/notification"
	pborder "github.com/thetherington/jobber-common/protogen/go/order"
	pbBuyer "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-order/internal/adapters/config"
	"github.com/thetherington/jobber-order/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

func (o *OrderService) CancelOrder(ctx context.Context, id string, paymentIntentId string, orderData *order.OrderMessage) (*order.OrderResponse, error) {
	if err := o.payment.RefundOrder(paymentIntentId); err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	var updatedOrder *order.OrderDocument

	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "cancelled", Value: true},
			{Key: "status", Value: "Cancelled"},
			{Key: "approvedAt", Value: time.Now()},
		},
	}}

	err := orderDB.FindOneAndUpdate(ctx, bson.M{"orderId": id}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// update seller orders from User microservice
	pbSeller := &pborder.SellerGigUpdate{
		Action:   pborder.Action_CancelOrder,
		SellerId: orderData.SellerId,
	}

	if data, err := proto.Marshal(pbSeller); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-seller-update", "user-seller", data); err != nil {
			slog.With("error", err).Error("CancelOrder: Failed to send offer to jobber-seller-update")
		}
	}

	// update buyer orders from User microservice
	pbBuyer := &pbBuyer.BuyerPayload{
		Action:        pbBuyer.Action_CANCELLED_GIG.Enum(),
		BuyerId:       utils.Ptr(orderData.BuyerId),
		PurchasedGigs: []string{orderData.PurchasedGigs},
	}

	if data, err := proto.Marshal(pbBuyer); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-buyer-update", "user-buyer", data); err != nil {
			slog.With("error", err).Error("CancelOrder: Failed to send offer to jobber-buyer-update")
		}
	}

	return &order.OrderResponse{Message: "Order cancelled successfully", Order: updatedOrder}, nil
}

func (o *OrderService) RequestExtension(ctx context.Context, id string, extension *order.ExtendedDelivery) (*order.OrderResponse, error) {
	// Validate Order payload
	if err := extension.Validate(validate); err != nil {
		slog.With("error", err).Debug("Request Extension Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	var updatedOrder *order.OrderDocument

	update := bson.D{{
		Key: "$set",
		Value: bson.D{{Key: "requestExtension", Value: bson.D{
			{Key: "originalDate", Value: extension.OriginalDate},
			{Key: "newDate", Value: extension.NewDate},
			{Key: "days", Value: extension.Days},
			{Key: "reason", Value: extension.Reason},
		}}},
	}}

	err := orderDB.FindOneAndUpdate(ctx, bson.M{"orderId": id}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// publish message send seller email
	pbEmail := &pbnotify.OrderEmailMessageDetails{
		ReceiverEmail:  &updatedOrder.BuyerEmail,
		BuyerUsername:  &updatedOrder.BuyerUsername,
		SellerUsername: &updatedOrder.SellerUsername,
		OriginalDate:   utils.Ptr(extension.OriginalDate.Format("2006-01-02 15:04:05")),
		NewDate:        utils.Ptr(extension.NewDate.Format("2006-01-02 15:04:05")),
		Reason:         &extension.Reason,
		OrderUrl:       utils.Ptr(fmt.Sprintf("%s/orders/%s/activities", config.Config.App.ClientUrl, updatedOrder.OrderId)),
		Template:       utils.Ptr("orderExtension"),
	}

	// send the order extension request information to the notification microservice to email the buyer
	if data, err := proto.Marshal(pbEmail); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-order-notification", "order-email", data); err != nil {
			slog.With("error", err).Error("RequestExtension: Failed to send offer to jobber-order-notification")
		}
	}

	return &order.OrderResponse{Message: "Order extenstion request", Order: updatedOrder}, nil
}

func (o *OrderService) DeliverOrder(ctx context.Context, id string, file *order.DeliveredWork) (*order.OrderResponse, error) {
	// Validate delivered work payload
	if err := file.Validate(validate); err != nil {
		slog.With("error", err).Debug("Deliver Order Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	randomCharacters := util.RandomString(20)

	url, err := o.image.UploadImage(ctx, file.File, randomCharacters, true, true)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	file.File = url

	var updatedOrder *order.OrderDocument

	update := bson.D{
		{
			Key:   "$push",
			Value: bson.D{{Key: "deliveredWork", Value: file}},
		},
		{
			Key: "$set",
			Value: bson.D{
				{Key: "delivered", Value: true},
				{Key: "status", Value: "Delivered"},
				{Key: "events.orderDelivered", Value: time.Now().Format(time.RFC3339)},
			},
		},
	}

	err = orderDB.FindOneAndUpdate(ctx, bson.M{"orderId": id}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// publish message send buyer email
	pbEmail := &pbnotify.OrderEmailMessageDetails{
		ReceiverEmail:  &updatedOrder.BuyerEmail,
		OrderId:        &updatedOrder.OrderId,
		BuyerUsername:  &updatedOrder.BuyerUsername,
		SellerUsername: &updatedOrder.SellerUsername,
		Title:          &updatedOrder.Offer.GigTitle,
		Description:    &updatedOrder.Offer.Description,
		OrderUrl:       utils.Ptr(fmt.Sprintf("%s/orders/%s/activities", config.Config.App.ClientUrl, updatedOrder.OrderId)),
		Template:       utils.Ptr("orderDelivered"),
	}

	// send the order delivery email to the notification microservice to email the buyer
	if data, err := proto.Marshal(pbEmail); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-order-notification", "order-email", data); err != nil {
			slog.With("error", err).Error("DeliverOrder: Failed to send offer to jobber-order-notification")
		}
	}

	return &order.OrderResponse{Message: "Order Delivered Received", Order: updatedOrder}, nil
}

func (o *OrderService) ApproveOrder(ctx context.Context, id string, orderData *order.OrderMessage) (*order.OrderResponse, error) {
	var updatedOrder *order.OrderDocument

	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "approved", Value: true},
			{Key: "status", Value: "Completed"},
			{Key: "approvedAt", Value: time.Now()},
		},
	}}

	err := orderDB.FindOneAndUpdate(ctx, bson.M{"orderId": id}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// update seller orders from User microservice
	pbSeller := &pborder.SellerGigUpdate{
		Action:   pborder.Action_ApproveOrder,
		SellerId: orderData.SellerId,
		OrderProps: &pborder.OrderProps{
			OngoingJobs:    &orderData.OngoingJobs,
			CompletedJobs:  &orderData.CompletedJobs,
			TotalEarnings:  &orderData.TotalEarnings,
			RecentDelivery: utils.CurrentDatetime(),
		},
	}

	if data, err := proto.Marshal(pbSeller); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-seller-update", "user-seller", data); err != nil {
			slog.With("error", err).Error("ApproveOrder: Failed to send offer to jobber-seller-update")
		}
	}

	return &order.OrderResponse{Message: "Order approved successfully", Order: updatedOrder}, nil
}

func (o *OrderService) RejectDeliveryDate(ctx context.Context, orderId string) (*order.OrderResponse, error) {
	var updatedOrder *order.OrderDocument

	update := bson.D{{
		Key:   "$unset",
		Value: bson.D{{Key: "requestExtension", Value: 1}},
	}}

	err := orderDB.FindOneAndUpdate(ctx, bson.M{"orderId": orderId}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// publish message send seller email
	pbEmail := &pbnotify.OrderEmailMessageDetails{
		ReceiverEmail:  &updatedOrder.SellerEmail,
		Subject:        utils.Ptr("Sorry: Your extension request was rejected"),
		BuyerUsername:  &updatedOrder.BuyerUsername,
		SellerUsername: &updatedOrder.SellerUsername,
		Header:         utils.Ptr("Request Rejected"),
		Message:        utils.Ptr("You can contact the buyer for more information."),
		Type:           utils.Ptr("rejected"),
		OrderUrl:       utils.Ptr(fmt.Sprintf("%s/orders/%s/activities", config.Config.App.ClientUrl, updatedOrder.OrderId)),
		Template:       utils.Ptr("orderExtensionApproval"),
	}

	// send the order extension rejection to the notification microservice to email the seller
	if data, err := proto.Marshal(pbEmail); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-order-notification", "order-email", data); err != nil {
			slog.With("error", err).Error("RejectDeliveryDate: Failed to send offer to jobber-order-notification")
		}
	}

	return &order.OrderResponse{Message: "Order extenstion rejected", Order: updatedOrder}, nil
}

func (o *OrderService) ApproveDeliveryDate(ctx context.Context, orderId string, data order.ExtendedDelivery) (*order.OrderResponse, error) {
	// Validate Extended Delivery payload
	if err := data.Validate(validate); err != nil {
		slog.With("error", err).Debug("Approve Delivery Date Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	var updatedOrder *order.OrderDocument

	update := bson.D{
		{
			Key:   "$unset",
			Value: bson.D{{Key: "requestExtension", Value: 1}},
		},
		{
			Key: "$set",
			Value: bson.D{
				{Key: "offer.deliveryInDays", Value: data.Days},
				{Key: "offer.newDeliveryDate", Value: data.NewDate.String()},
				{Key: "offer.reason", Value: data.Reason},
				{Key: "events.deliveryDateUpdate", Value: data.DeliveryDateUpdate},
			},
		},
	}

	err := orderDB.FindOneAndUpdate(ctx, bson.M{"orderId": orderId}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// publish message send seller email
	pbEmail := &pbnotify.OrderEmailMessageDetails{
		ReceiverEmail:  &updatedOrder.SellerEmail,
		Subject:        utils.Ptr("Congratulations: Your extension request was approved"),
		BuyerUsername:  &updatedOrder.BuyerUsername,
		SellerUsername: &updatedOrder.SellerUsername,
		Header:         utils.Ptr("Request Accepted"),
		Message:        utils.Ptr("You can continue working on the order."),
		Type:           utils.Ptr("accepted"),
		OrderUrl:       utils.Ptr(fmt.Sprintf("%s/orders/%s/activities", config.Config.App.ClientUrl, updatedOrder.OrderId)),
		Template:       utils.Ptr("orderExtensionApproval"),
	}

	// send the order extension rejection to the notification microservice to email the seller
	if data, err := proto.Marshal(pbEmail); err == nil {
		if err := o.queue.PublishDirectMessage("jobber-order-notification", "order-email", data); err != nil {
			slog.With("error", err).Error("ApproveDeliveryDate: Failed to send offer to jobber-order-notification")
		}
	}

	return &order.OrderResponse{Message: "Order extenstion approved", Order: updatedOrder}, nil
}

func (o *OrderService) UpdateOrderReview(ctx context.Context, data *review.ReviewMessageDetails) (*order.OrderResponse, string, error) {
	var updatedOrder *order.OrderDocument

	reviewType := utils.FirstToLower(data.Type)

	update := bson.D{{
		Key: "$set", Value: bson.D{
			{Key: reviewType, Value: bson.D{
				{Key: "rating", Value: data.Rating},
				{Key: "review", Value: data.Review},
				{Key: "date", Value: data.CreatedAt.Format("2006-01-02 15:04:05")},
			},
			},
			{Key: fmt.Sprintf("events.%s", reviewType), Value: data.CreatedAt.Format("2006-01-02 15:04:05")},
		},
	}}

	err := orderDB.FindOneAndUpdate(ctx, bson.M{"orderId": data.OrderId}, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedOrder)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", svc.NewError(svc.ErrNotFound, fmt.Errorf("order does not exist"))
		}

		return nil, "", svc.NewError(svc.ErrInternalFailure, err)
	}

	userTo := updatedOrder.BuyerUsername
	if data.Type == "BuyerReview" {
		userTo = updatedOrder.SellerUsername
	}

	msg := fmt.Sprintf("%s left you a %d star review", userTo, data.Rating)

	return &order.OrderResponse{Message: msg, Order: updatedOrder}, userTo, nil
}
