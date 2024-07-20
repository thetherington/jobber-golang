package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thetherington/jobber-common/models/order"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

// OrderHandler represents the HTTP handler for order service requests
type OrderHandler struct {
	svc port.OrderService
}

// /api/gateway/v1/order
func (oh OrderHandler) Routes(router chi.Router) {
	router.Get("/notification/{userTo}", oh.GetNotificationsById)
	router.Get("/{orderId}", oh.GetOrderById)
	router.Get("/seller/{sellerId}", oh.GetSellerOrders)
	router.Get("/buyer/{buyerId}", oh.GetBuyerOrders)

	router.Post("/", oh.CreateOrder)
	router.Post("/create-payment-intent", oh.CreatePaymentIntent)

	router.Put("/cancel/{orderId}", oh.CancelOrder)
	router.Put("/extension/{orderId}", oh.RequestExtension)
	router.Put("/deliver-order/{orderId}", oh.DeliverOrder)
	router.Put("/approve-order/{orderId}", oh.ApproveOrder)
	router.Put("/gig/{type}/{orderId}", oh.DeliveryDate)

	router.Put("/notification/mark-as-read", oh.MarkNotificationAsRead)
}

// NewOrderHandler creates a new OrderHandler instance
func NewOrderHandler(svc port.OrderService) *OrderHandler {
	return &OrderHandler{
		svc,
	}
}

func (oh *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order *order.OrderDocument

	// unmarshal the request body into the order
	if err := ReadJSON(w, r, &order); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the order microservice via grpc client
	resp, err := oh.svc.CreateOrder(r.Context(), order)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusCreated, &resp); err != nil {
		slog.With("error", err).Error("CreateOrder: failed to write http response")
	}
}

func (oh *OrderHandler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	type PaymentIntent struct {
		Price   float32 `json:"price"`
		BuyerId string  `json:"buyerId"`
	}

	var payload *PaymentIntent

	// unmarshal the request body into the order
	if err := ReadJSON(w, r, &payload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the order microservice via grpc client
	resp, err := oh.svc.CreatePaymentIntent(r.Context(), payload.Price, payload.BuyerId)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusCreated, &resp); err != nil {
		slog.With("error", err).Error("CreatePaymentIntent: failed to write http response")
	}
}

func (oh *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	// send a request to the order microservice via grpc client
	resp, err := oh.svc.GetOrderById(r.Context(), chi.URLParam(r, "orderId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetSellerGigs: failed to write http response")
	}
}

func (oh *OrderHandler) GetSellerOrders(w http.ResponseWriter, r *http.Request) {
	// send a request to the order microservice via grpc client
	resp, err := oh.svc.GetSellerOrders(r.Context(), chi.URLParam(r, "sellerId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetSellerGigs: failed to write http response")
	}
}

func (oh *OrderHandler) GetBuyerOrders(w http.ResponseWriter, r *http.Request) {
	// send a request to the order microservice via grpc client
	resp, err := oh.svc.GetBuyerOrders(r.Context(), chi.URLParam(r, "buyerId"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetBuyerOrders: failed to write http response")
	}
}

func (oh *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	type Cancellation struct {
		PaymentIntentId string              `json:"paymentIntentId"`
		OrderData       *order.OrderMessage `json:"orderData"`
	}

	var payload *Cancellation

	// unmarshal the request body into the cancellation payload
	if err := ReadJSON(w, r, &payload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the order microservice via grpc client
	msg, err := oh.svc.CancelOrder(r.Context(), chi.URLParam(r, "orderId"), payload.PaymentIntentId, payload.OrderData)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("CancelOrder: failed to write http response")
	}

}

func (oh *OrderHandler) RequestExtension(w http.ResponseWriter, r *http.Request) {
	var payload *order.ExtendedDelivery

	// unmarshal the request body into the cancellation payload
	if err := ReadJSON(w, r, &payload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := oh.svc.RequestExtension(r.Context(), chi.URLParam(r, "orderId"), payload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("RequestExtension: failed to write http response")
	}
}

func (oh *OrderHandler) DeliverOrder(w http.ResponseWriter, r *http.Request) {
	var payload *order.DeliveredWork

	// unmarshal the request body into the cancellation payload
	if err := ReadJSON(w, r, &payload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := oh.svc.DeliverOrder(r.Context(), chi.URLParam(r, "orderId"), payload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("RequestExtension: failed to write http response")
	}
}

func (oh *OrderHandler) ApproveOrder(w http.ResponseWriter, r *http.Request) {
	var payload *order.OrderMessage

	// unmarshal the request body into the cancellation payload
	if err := ReadJSON(w, r, &payload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := oh.svc.ApproveOrder(r.Context(), chi.URLParam(r, "orderId"), payload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("ApproveOrder: failed to write http response")
	}
}

func (oh *OrderHandler) DeliveryDate(w http.ResponseWriter, r *http.Request) {
	var payload *order.ExtendedDelivery

	// unmarshal the request body into the cancellation payload
	if err := ReadJSON(w, r, &payload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := oh.svc.DeliveryDate(r.Context(), chi.URLParam(r, "orderId"), chi.URLParam(r, "type"), payload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("DeliveryDate: failed to write http response")
	}
}

func (oh *OrderHandler) GetNotificationsById(w http.ResponseWriter, r *http.Request) {
	resp, err := oh.svc.GetNotificationsById(r.Context(), chi.URLParam(r, "userTo"))
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetNotificationsById: failed to write http response")
	}
}

func (oh *OrderHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		NotificationId string `json:"notificationId"`
	}

	var payload *Payload

	// unmarshal the request body
	if err := ReadJSON(w, r, &payload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp, err := oh.svc.MarkNotificationAsRead(r.Context(), payload.NotificationId)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("MarkNotificationAsRead: failed to write http response")
	}
}
