package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thetherington/jobber-common/models/chat"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

// ChatHandler represents the HTTP handler for gig service requests
type ChatHandler struct {
	svc port.ChatService
}

// /api/gateway/v1/message
func (ch ChatHandler) Routes(router chi.Router) {
	router.Get("/conversation/{senderUsername}/{receiverUsername}", ch.GetConversation)
	router.Get("/conversations/{username}", ch.GetConversationList)
	router.Get("/{senderUsername}/{receiverUsername}", ch.GetMessages)
	router.Get("/{conversationId}", ch.GetUserMessages)

	router.Post("/", ch.CreateMessage)

	router.Put("/offer", ch.UpdateOffer)
	router.Put("/mark-as-read", ch.MarkAsRead)
	router.Put("/mark-multiple-as-read", ch.MarkMultipleAsRead)
}

// NewChatHandler creates a new ChatHandler instance
func NewChatHandler(svc port.ChatService) *ChatHandler {
	return &ChatHandler{
		svc,
	}
}

func (ch *ChatHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var messagePayload *chat.MessageDocument

	// unmarshal the request body into the messagePayload
	if err := ReadJSON(w, r, &messagePayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := ch.svc.AddMessage(r.Context(), messagePayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusCreated, &resp); err != nil {
		slog.With("error", err).Error("CreateMessage: failed to write http response")
	}
}

func (ch *ChatHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	senderUsername := chi.URLParam(r, "senderUsername")
	receiverUsername := chi.URLParam(r, "receiverUsername")

	resp, err := ch.svc.GetConversation(r.Context(), senderUsername, receiverUsername)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetConversation: failed to write http response")
	}
}

func (ch *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	senderUsername := chi.URLParam(r, "senderUsername")
	receiverUsername := chi.URLParam(r, "receiverUsername")

	resp, err := ch.svc.GetMessages(r.Context(), senderUsername, receiverUsername)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetMessages: failed to write http response")
	}
}

func (ch *ChatHandler) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	resp, err := ch.svc.GetUserMessages(r.Context(), chi.URLParam(r, "conversationId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetMessages: failed to write http response")
	}
}

func (ch *ChatHandler) GetConversationList(w http.ResponseWriter, r *http.Request) {
	resp, err := ch.svc.GetUserConversationList(r.Context(), chi.URLParam(r, "username"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetMessages: failed to write http response")
	}
}

func (ch *ChatHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MessageId string `json:"messageId"`
		Type      string `json:"type"`
	}

	// unmarshal the request body into the OfferUpdate
	if err := ReadJSON(w, r, &body); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := ch.svc.UpdateOffer(r.Context(), body.MessageId, body.Type)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("UpdateOffer: failed to write http response")
	}
}

func (ch *ChatHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MessageId string `json:"messageId"`
	}

	// unmarshal the request body into the MarkAsRead
	if err := ReadJSON(w, r, &body); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := ch.svc.MarkMessageAsRead(r.Context(), body.MessageId)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("MarkAsRead: failed to write http response")
	}
}

func (ch *ChatHandler) MarkMultipleAsRead(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MessageId        string `json:"messageId"`
		ReceiverUsername string `json:"receiverUsername"`
		SenderUsername   string `json:"senderUsername"`
	}

	// unmarshal the request body into the MarkMultipleAsRead
	if err := ReadJSON(w, r, &body); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := ch.svc.MarkManyMessagesAsRead(r.Context(), body.SenderUsername, body.ReceiverUsername, body.MessageId)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	msg := &struct {
		Message string `json:"message"`
	}{
		Message: resp,
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &msg); err != nil {
		slog.With("error", err).Error("MarkMultipleAsRead: failed to write http response")
	}
}
