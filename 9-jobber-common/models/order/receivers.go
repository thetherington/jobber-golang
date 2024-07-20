package order

import (
	"github.com/go-playground/validator/v10"
	pb "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-common/utils"
)

func (s *OrderDocument) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[OrderDocument](*s, validate)
}

func (s *ExtendedDelivery) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[ExtendedDelivery](*s, validate)
}

func (s *Offer) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[Offer](*s, validate)
}

func (s *DeliveredWork) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[DeliveredWork](*s, validate)
}

func (o *OrderDocument) MarshalToProto() *pb.OrderDocument {
	pbmsg := &pb.OrderDocument{
		OrderId:             o.OrderId,
		InvoiceId:           o.InvoiceId,
		PaymentIntent:       o.PaymentIntent,
		GigId:               o.GigId,
		SellerId:            o.SellerId,
		SellerUsername:      o.SellerUsername,
		SellerImage:         o.SellerImage,
		SellerEmail:         o.SellerEmail,
		GigCoverImage:       o.GigCoverImage,
		GigMainTitle:        o.GigMainTitle,
		GigBasicTitle:       o.GigBasicTitle,
		GigBasicDescription: o.GigBasicDescription,
		BuyerId:             o.BuyerId,
		BuyerUsername:       o.BuyerUsername,
		BuyerEmail:          o.BuyerEmail,
		BuyerImage:          o.BuyerImage,
		Status:              o.Status,
		Requirements:        o.Requirements,
		Quantity:            o.Quantity,
		Price:               o.Price,
		ServiceFee:          o.ServiceFee,
		Approved:            o.Approved,
		Delivered:           o.Delivered,
		Cancelled:           o.Cancelled,
		ApprovedAt:          utils.ToDateTimeOrNil(o.ApprovedAt),
		DateOrdered:         utils.ToDateTime(o.DateOrdered),
		Offer: &pb.Offer{
			GigTitle:        o.Offer.GigTitle,
			Price:           o.Offer.Price,
			Description:     o.Offer.Description,
			DeliverInDays:   o.Offer.DeliveryInDays,
			OldDeliveryDate: o.Offer.OldDeliveryDate,
			NewDeliveryDate: o.Offer.NewDeliveryDate,
			Accepted:        o.Offer.Accepted,
			Cancelled:       o.Offer.Cancelled,
			Reason:          o.Offer.Reason,
		},
		DeliveredWork:    make([]*pb.DeliveredWork, 0),
		Events:           nil,
		BuyerReview:      nil,
		SellerReview:     nil,
		RequestExtension: nil,
	}

	if o.DeliveredWork != nil {
		var deliveredWork []*pb.DeliveredWork

		for _, m := range o.DeliveredWork {
			deliveredWork = append(deliveredWork, &pb.DeliveredWork{
				Message:  m.Message,
				File:     m.File,
				FileType: m.FileType,
				FileSize: m.FileSize,
				FileName: m.FileName,
			})
		}

		pbmsg.DeliveredWork = deliveredWork
	}

	if o.Events != nil {
		pbmsg.Events = &pb.OrderEvents{
			PlaceOrder:         o.Events.PlaceOrder,
			Requirements:       o.Events.Requirements,
			OrderStarted:       o.Events.OrderStarted,
			DeliveryDateUpdate: o.Events.DeliveryDateUpdate,
			OrderDelivered:     o.Events.OrderStarted,
			BuyerReview:        o.Events.BuyerReview,
			SellerReview:       o.Events.SellerReview,
		}
	}

	if o.BuyerReview != nil {
		pbmsg.BuyerReview = &pb.OrderReview{
			Rating: o.BuyerReview.Rating,
			Review: o.BuyerReview.Review,
			Date:   o.BuyerReview.Date,
		}
	}

	if o.SellerReview != nil {
		pbmsg.SellerReview = &pb.OrderReview{
			Rating: o.SellerReview.Rating,
			Review: o.SellerReview.Review,
			Date:   o.SellerReview.Date,
		}
	}

	if o.RequestExtension != nil {
		pbmsg.RequestExtension = &pb.ExtendedDelivery{
			OriginalDate:        utils.ToDateTimeOrNil(o.RequestExtension.OriginalDate),
			NewDate:             utils.ToDateTimeOrNil(o.RequestExtension.NewDate),
			Days:                o.RequestExtension.Days,
			Reason:              o.RequestExtension.Reason,
			DeliveryDateUpdated: o.RequestExtension.DeliveryDateUpdate,
		}
	}

	return pbmsg
}

func (o *OrderMessage) MarshalToProto() *pb.OrderMessage {
	return &pb.OrderMessage{
		SellerId:       ReturnPtrOrNil(o.SellerId),
		BuyerId:        ReturnPtrOrNil(o.BuyerId),
		OngoingJobs:    &o.OngoingJobs,
		CompletedJobs:  &o.CompletedJobs,
		TotalEarnings:  &o.TotalEarnings,
		PurchasedGigs:  ReturnPtrOrNil(o.PurchasedGigs),
		RecentDelivery: ReturnPtrOrNil(o.RecentDelivery),
		Type:           ReturnPtrOrNil(o.Type),
		ReceiverEmail:  ReturnPtrOrNil(o.ReceiverEmail),
		Username:       ReturnPtrOrNil(o.Username),
		Template:       ReturnPtrOrNil(o.Template),
		Sender:         ReturnPtrOrNil(o.Sender),
		OfferLink:      ReturnPtrOrNil(o.OfferLink),
		Amount:         ReturnPtrOrNil(o.Amount),
		BuyerUsername:  ReturnPtrOrNil(o.BuyerUsername),
		SellerUsername: ReturnPtrOrNil(o.SellerUsername),
		Title:          ReturnPtrOrNil(o.Title),
		Description:    ReturnPtrOrNil(o.Description),
		DeliveryDays:   ReturnPtrOrNil(o.DeliveryDays),
		OrderId:        ReturnPtrOrNil(o.OrderId),
		InvoiceId:      ReturnPtrOrNil(o.InvoiceId),
		OrderDue:       ReturnPtrOrNil(o.OrderDue),
		Requirements:   ReturnPtrOrNil(o.Requirements),
		OrderUrl:       ReturnPtrOrNil(o.OrderUrl),
		OriginalDate:   ReturnPtrOrNil(o.OriginalDate),
		NewDate:        ReturnPtrOrNil(o.NewDate),
		Reason:         ReturnPtrOrNil(o.Reason),
		Subject:        ReturnPtrOrNil(o.Subject),
		Header:         ReturnPtrOrNil(o.Header),
		Total:          ReturnPtrOrNil(o.Total),
		Message:        ReturnPtrOrNil(o.Message),
		ServiceFee:     ReturnPtrOrNil(o.ServiceFee),
	}
}

func ReturnPtrOrNil(s string) *string {
	if s != "" {
		return &s
	}

	return nil
}

func CreateOrderMessage(m *pb.OrderMessage) *OrderMessage {
	return &OrderMessage{
		SellerId:       m.GetSellerId(),
		BuyerId:        m.GetBuyerId(),
		OngoingJobs:    m.GetOngoingJobs(),
		CompletedJobs:  m.GetCompletedJobs(),
		TotalEarnings:  m.GetTotalEarnings(),
		PurchasedGigs:  m.GetPurchasedGigs(),
		RecentDelivery: m.GetRecentDelivery(),
		Type:           m.GetType(),
		ReceiverEmail:  m.GetReceiverEmail(),
		Username:       m.GetUsername(),
		Template:       m.GetTemplate(),
		Sender:         m.GetSender(),
		OfferLink:      m.GetOfferLink(),
		Amount:         m.GetAmount(),
		BuyerUsername:  m.GetBuyerUsername(),
		SellerUsername: m.GetSellerUsername(),
		Title:          m.GetTitle(),
		Description:    m.GetDescription(),
		DeliveryDays:   m.GetDeliveryDays(),
		OrderId:        m.GetOrderId(),
		InvoiceId:      m.GetInvoiceId(),
		OrderDue:       m.GetOrderDue(),
		Requirements:   m.GetRequirements(),
		OrderUrl:       m.GetOrderUrl(),
		OriginalDate:   m.GetOriginalDate(),
		NewDate:        m.GetNewDate(),
		Reason:         m.GetReason(),
		Subject:        m.GetSubject(),
		Header:         m.GetHeader(),
		Total:          m.GetTotal(),
		Message:        m.GetMessage(),
		ServiceFee:     m.GetServiceFee(),
	}
}

func CreateOrderDocument(m *pb.OrderDocument) *OrderDocument {
	o := &OrderDocument{
		OrderId:             m.GetOrderId(),
		InvoiceId:           m.GetInvoiceId(),
		PaymentIntent:       m.GetPaymentIntent(),
		GigId:               m.GetGigId(),
		SellerId:            m.GetSellerId(),
		SellerUsername:      m.GetSellerUsername(),
		SellerImage:         m.GetSellerImage(),
		SellerEmail:         m.GetSellerEmail(),
		GigCoverImage:       m.GetGigCoverImage(),
		GigMainTitle:        m.GetGigMainTitle(),
		GigBasicTitle:       m.GetGigBasicTitle(),
		GigBasicDescription: m.GetGigBasicDescription(),
		BuyerId:             m.GetBuyerId(),
		BuyerUsername:       m.GetBuyerUsername(),
		BuyerEmail:          m.GetBuyerEmail(),
		BuyerImage:          m.GetBuyerImage(),
		Status:              m.GetStatus(),
		Requirements:        m.GetRequirements(),
		Quantity:            m.GetQuantity(),
		Price:               m.GetPrice(),
		ServiceFee:          m.GetServiceFee(),
		Approved:            m.GetApproved(),
		Delivered:           m.GetDelivered(),
		Cancelled:           m.GetCancelled(),
		ApprovedAt:          utils.ToTimeOrNil(m.ApprovedAt),
		DateOrdered:         utils.ToTimeOrNil(m.DateOrdered),
		Offer: &Offer{
			GigTitle:        m.Offer.GetGigTitle(),
			Price:           m.Offer.GetPrice(),
			Description:     m.Offer.GetDescription(),
			DeliveryInDays:  m.Offer.GetDeliverInDays(),
			OldDeliveryDate: m.Offer.GetOldDeliveryDate(),
			NewDeliveryDate: m.Offer.GetNewDeliveryDate(),
			Accepted:        m.Offer.GetAccepted(),
			Cancelled:       m.Offer.GetCancelled(),
			Reason:          m.Offer.GetReason(),
		},
		DeliveredWork:    make([]*DeliveredWork, 0),
		Events:           nil,
		BuyerReview:      nil,
		SellerReview:     nil,
		RequestExtension: nil,
	}

	if m.DeliveredWork != nil {
		var deliveredWork []*DeliveredWork

		for _, w := range m.DeliveredWork {
			deliveredWork = append(deliveredWork, &DeliveredWork{
				Message:  w.GetMessage(),
				File:     w.GetFile(),
				FileType: w.GetFileType(),
				FileSize: w.GetFileSize(),
				FileName: w.GetFileName(),
			})
		}

		o.DeliveredWork = deliveredWork
	}

	if m.Events != nil {
		o.Events = &OrderEvents{
			PlaceOrder:         m.Events.GetPlaceOrder(),
			Requirements:       m.Events.GetRequirements(),
			OrderStarted:       m.Events.GetOrderStarted(),
			DeliveryDateUpdate: m.Events.GetDeliveryDateUpdate(),
			OrderDelivered:     m.Events.GetOrderDelivered(),
			BuyerReview:        m.Events.GetBuyerReview(),
			SellerReview:       m.Events.GetSellerReview(),
		}
	}

	if m.BuyerReview != nil {
		o.BuyerReview = &OrderReview{
			Rating: m.BuyerReview.GetRating(),
			Review: m.BuyerReview.GetReview(),
			Date:   m.BuyerReview.GetDate(),
		}
	}

	if m.SellerReview != nil {
		o.SellerReview = &OrderReview{
			Rating: m.SellerReview.GetRating(),
			Review: m.SellerReview.GetReview(),
			Date:   m.SellerReview.GetDate(),
		}
	}

	if m.RequestExtension != nil {
		o.RequestExtension = &ExtendedDelivery{
			OriginalDate:       utils.ToTimeOrNil(m.RequestExtension.OriginalDate),
			NewDate:            utils.ToTimeOrNil(m.RequestExtension.NewDate),
			Days:               m.RequestExtension.GetDays(),
			Reason:             m.RequestExtension.GetReason(),
			DeliveryDateUpdate: m.RequestExtension.GetDeliveryDateUpdated(),
		}
	}

	return o
}
