package order

import "time"

type Notification struct {
	Id               string     `json:"_id"              bson:"_id,omitempty"`
	UserTo           string     `json:"userTo"           bson:"userTo"`
	SenderUsername   string     `json:"senderUsername"   bson:"senderUsername"`
	SenderPicture    string     `json:"senderPicture"    bson:"senderPicture"`
	ReceiverUsername string     `json:"receiverUsername" bson:"receiverUsername"`
	ReceiverPicture  string     `json:"receiverPicture"  bson:"receiverPicture"`
	IsRead           bool       `json:"isRead"           bson:"isRead"`
	Message          string     `json:"message"          bson:"message"`
	OrderId          string     `json:"orderId"          bson:"orderId"`
	CreatedAt        *time.Time `json:"createdAt"        bson:"createdAt"`
}

type NotificationResponse struct {
	Message      string        `json:"message"`
	Notification *Notification `json:"notification"`
}

type NotificationsResponse struct {
	Message       string          `json:"message"`
	Notifications []*Notification `json:"notifications"`
}
