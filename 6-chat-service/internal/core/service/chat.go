package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/thetherington/jobber-chat/internal/adapters/config"
	"github.com/thetherington/jobber-chat/internal/core/port"
	"github.com/thetherington/jobber-chat/internal/core/util"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/models/chat"
	pb "github.com/thetherington/jobber-common/protogen/go/notification"
	"github.com/thetherington/jobber-common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

var (
	messageDB      *mongo.Collection
	conversationDB *mongo.Collection
	validate       *validator.Validate
)

/**
 * ChatService implements
 */
type ChatService struct {
	queue port.ChatProducer
	image port.ImageUploader
}

// NewChatService creates a new chat service instance
func NewChatService(db *mongo.Database, queue port.ChatProducer, image port.ImageUploader) *ChatService {
	validate = validator.New(validator.WithRequiredStructEnabled())

	conversationDB = db.Collection("Conversation")
	messageDB = db.Collection("Message")

	conversationDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "conversationId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "senderUsername", Value: 1}}},
		{Keys: bson.D{{Key: "receiverUsername", Value: 1}}},
	})

	messageDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "conversationId", Value: 1}}},
		{Keys: bson.D{{Key: "senderUsername", Value: 1}}},
		{Keys: bson.D{{Key: "receiverUsername", Value: 1}}},
	})

	return &ChatService{
		queue: queue,
		image: image,
	}
}

func (c *ChatService) CreateConversation(ctx context.Context, conversationId, sender, receiver string) error {
	conversation := chat.ConversationDocument{
		CoversationId:    conversationId,
		SenderUsername:   sender,
		ReceiverUsername: receiver,
	}

	if _, err := conversationDB.InsertOne(ctx, conversation); err != nil {
		slog.With("error", err).Error("failed to create conversation")
		return svc.NewError(svc.ErrInternalFailure, err)
	}

	return nil
}

func (c *ChatService) AddMessage(ctx context.Context, msg *chat.MessageDocument) (*chat.MessageResponse, error) {
	// Validate Message payload
	if err := msg.Validate(validate); err != nil {
		slog.With("error", err).Debug("AddMessage Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// Cloudinary image upload
	if msg.File != "" {
		randomCharacters := util.RandomString(20)

		url, err := c.image.UploadImage(ctx, msg.File, randomCharacters, true, true)
		if err != nil {
			return nil, svc.NewError(svc.ErrInternalFailure, err)
		}

		msg.File = url
	}

	// Dispatch to notification service via the users microservice as relay if offer in message
	if msg.Offer != nil {
		pbMsg := &pb.OrderEmailMessageDetails{
			Sender:         &msg.SenderUsername,
			Amount:         utils.Ptr(strconv.FormatFloat(float64(msg.Offer.Price), 'f', 2, 64)),
			BuyerUsername:  &msg.ReceiverUsername,
			SellerUsername: &msg.SenderUsername,
			Title:          &msg.Offer.GigTitle,
			Description:    &msg.Offer.Description,
			OfferLink:      utils.Ptr(fmt.Sprintf("%s/inbox/%s/%s", config.Config.App.ClientUrl, msg.SenderUsername, msg.CoversationId)),
			DeliveryDays:   utils.Ptr(strconv.Itoa(int(msg.Offer.DeliveryInDays))),
			Template:       utils.Ptr("offer"),
		}

		// send the offer information to the users microservice to pickup the buyer email and sent to the notification service
		if data, err := proto.Marshal(pbMsg); err == nil {
			if err := c.queue.PublishDirectMessage("jobber-relay-notification", "relay-notification", data); err != nil {
				slog.With("error", err).Error("AddMessage: Failed to send offer to jobber-relay-notification")
			}
		}
	}

	// Create a conversation record if the client says there's no conversation
	if !msg.HasConversationId {
		c.CreateConversation(ctx, msg.CoversationId, msg.SenderUsername, msg.ReceiverUsername)
	}

	resp, err := messageDB.InsertOne(ctx, msg)
	if err != nil {
		slog.With("error", err).Error("failed to save message into db")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	objId, ok := resp.InsertedID.(primitive.ObjectID)
	if ok {
		msg.Id = objId.Hex()
	}

	return &chat.MessageResponse{Message: "Message added", ConversationId: msg.CoversationId, MessageData: msg}, nil
}
