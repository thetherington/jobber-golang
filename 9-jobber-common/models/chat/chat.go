package chat

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/thetherington/jobber-common/models/order"
	pb "github.com/thetherington/jobber-common/protogen/go/chat"
	pbOrder "github.com/thetherington/jobber-common/protogen/go/order"
	"github.com/thetherington/jobber-common/utils"
)

type ConversationDocument struct {
	Id               string `json:"_id,omitempty"    bson:"_id,omitempty"`
	CoversationId    string `json:"conversationId"   bson:"conversationId"`
	SenderUsername   string `json:"senderUsername"   bson:"senderUsername"`
	ReceiverUsername string `json:"receiverUsername" bson:"receiverUsername"`
}

type MessageDocument struct {
	Id                string       `json:"_id"                bson:"_id,omitempty"      `
	CoversationId     string       `json:"conversationId"     bson:"conversationId"     `
	Body              string       `json:"body"               bson:"body"               `
	Url               string       `json:"url"                bson:"url"                `
	File              string       `json:"file"               bson:"file"               `
	FileType          string       `json:"fileType"           bson:"fileType"           `
	FileSize          string       `json:"fileSize"           bson:"fileSize"           `
	FileName          string       `json:"fileName"           bson:"fileName"           `
	GigId             string       `json:"gigId"              bson:"gigId"              `
	SellerId          string       `json:"sellerId"           bson:"sellerId"           validate:"required"   errmsg:"Please provide the seller id"`
	BuyerId           string       `json:"buyerId"            bson:"buyerId"            validate:"required"   errmsg:"Please provide the buyer id"`
	SenderUsername    string       `json:"senderUsername"     bson:"senderUsername"     validate:"required"   errmsg:"Please provide the sender username"`
	SenderPicture     string       `json:"senderPicture"      bson:"senderPicture"      validate:"required"   errmsg:"Please provide the sender picture"`
	ReceiverUsername  string       `json:"receiverUsername"   bson:"receiverUsername"   validate:"required"   errmsg:"Please provide the receiver username"`
	ReceiverPicture   string       `json:"receiverPicture"    bson:"receiverPicture"    validate:"required"   errmsg:"Please provide the receiver picture"`
	IsRead            bool         `json:"isRead"             bson:"isRead"             `
	HasOffer          bool         `json:"hasOffer"           bson:"hasOffer"           `
	Offer             *order.Offer `json:"offer,omitempty"    bson:"offer,omitempty"    `
	HasConversationId bool         `json:"hasConversationId"  bson:"hasConversationId"  `
	CreatedAt         *time.Time   `json:"createdAt"          bson:"createdAt"          `
}

func (s *MessageDocument) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[MessageDocument](*s, validate)
}

type MessageDetails struct {
	ReceiverEmail  string `json:"receiverEmail"`
	Sender         string `json:"sender"`
	OfferLink      string `json:"offerLink"`
	Amount         string `json:"amount"`
	BuyerUsername  string `json:"buyerUsername"`
	SellerUsername string `json:"sellerUsername"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	DeliveryInDays string `json:"deliveryInDays"`
	Template       string `json:"template"`
}

type MessagesResponse struct {
	Message  string             `json:"message"`
	Messages []*MessageDocument `json:"messages"`
}

type ConversationsResponse struct {
	Message       string                  `json:"message"`
	Conversations []*ConversationDocument `json:"conversations"`
}

type ConversationListResponse struct {
	Message       string             `json:"message"`
	Conversations []*MessageDocument `json:"conversations"`
}

type MessageResponse struct {
	Message        string           `json:"message"`
	ConversationId string           `json:"conversationId,omitempty"`
	MessageData    *MessageDocument `json:"messageData"`
}

type MarkSingleMessageResponse struct {
	Message       string           `json:"message"`
	SingleMessage *MessageDocument `json:"singleMessage"`
}

func CreateMessageDocument(m *pb.MessageDocument) *MessageDocument {
	msg := &MessageDocument{
		Id:                m.Id,
		CoversationId:     m.CoversationId,
		Body:              m.Body,
		Url:               m.Url,
		File:              m.File,
		FileType:          m.FileType,
		FileSize:          m.FileSize,
		FileName:          m.FileName,
		GigId:             m.GigId,
		SellerId:          m.SellerId,
		BuyerId:           m.BuyerId,
		SenderUsername:    m.SenderUsername,
		SenderPicture:     m.SenderPicture,
		ReceiverUsername:  m.ReceiverUsername,
		ReceiverPicture:   m.ReceiverPicture,
		IsRead:            m.IsRead,
		HasOffer:          m.HasOffer,
		HasConversationId: m.HasConversationId,
		Offer:             nil,
		CreatedAt:         utils.ToTime(m.CreatedAt),
	}

	if m.HasOffer && m.Offer != nil {
		offer := &order.Offer{
			GigTitle:        m.Offer.GigTitle,
			Price:           m.Offer.Price,
			Description:     m.Offer.Description,
			DeliveryInDays:  m.Offer.DeliverInDays,
			OldDeliveryDate: m.Offer.OldDeliveryDate,
			NewDeliveryDate: m.Offer.NewDeliveryDate,
			Accepted:        m.Offer.Accepted,
			Cancelled:       m.Offer.Cancelled,
			Reason:          m.Offer.Reason,
		}

		msg.Offer = offer
	}

	return msg
}

func CreateProtoMessageDocument(msg *MessageDocument) *pb.MessageDocument {
	req := &pb.MessageDocument{
		Id:                msg.Id,
		CoversationId:     msg.CoversationId,
		Body:              msg.Body,
		Url:               msg.Url,
		File:              msg.File,
		FileType:          msg.FileType,
		FileSize:          msg.FileSize,
		FileName:          msg.FileName,
		GigId:             msg.GigId,
		SellerId:          msg.SellerId,
		BuyerId:           msg.BuyerId,
		SenderUsername:    msg.SenderUsername,
		SenderPicture:     msg.SenderPicture,
		ReceiverUsername:  msg.ReceiverUsername,
		ReceiverPicture:   msg.ReceiverPicture,
		IsRead:            msg.IsRead,
		HasOffer:          msg.HasOffer,
		HasConversationId: msg.HasConversationId,
		CreatedAt:         utils.ToDateTime(msg.CreatedAt),
	}

	if msg.HasOffer && msg.Offer != nil {
		offer := &pbOrder.Offer{
			GigTitle:        msg.Offer.GigTitle,
			Price:           msg.Offer.Price,
			Description:     msg.Offer.Description,
			DeliverInDays:   msg.Offer.DeliveryInDays,
			OldDeliveryDate: msg.Offer.OldDeliveryDate,
			NewDeliveryDate: msg.Offer.NewDeliveryDate,
			Accepted:        msg.Offer.Accepted,
			Cancelled:       msg.Offer.Cancelled,
			Reason:          msg.Offer.Reason,
		}

		req.Offer = offer
	}

	return req
}
