package order

import "time"

type Offer struct {
	GigTitle        string  `json:"gigTitle"         bson:"gigTitle"         validate:"required"  errmsg:"Please provide the gig title"`
	Price           float32 `json:"price"            bson:"price"            validate:"required"  errmsg:"Please provide the price of the offer"`
	Description     string  `json:"description"      bson:"description"      validate:"required"  errmsg:"Please provide the offer description"`
	DeliveryInDays  int32   `json:"deliveryInDays"   bson:"deliveryInDays"   validate:"required"  errmsg:"Please provide delivery number in days"`
	OldDeliveryDate string  `json:"oldDeliveryDate"  bson:"oldDeliveryDate"  validate:"required"  errmsg:"Please provide the old delivery date"`
	NewDeliveryDate string  `json:"newDeliveryDate"  bson:"newDeliveryDate"  `
	Accepted        bool    `json:"accepted"         bson:"accepted"         `
	Cancelled       bool    `json:"cancelled"        bson:"cancelled"        `
	Reason          string  `json:"reason"           bson:"reason"           `
}

type OrderDocument struct {
	OrderId             string            `json:"orderId"                     bson:"orderId"               validate:"required"        errmsg:"Please provide the order id"`
	InvoiceId           string            `json:"invoiceId"                   bson:"invoiceId"             validate:"required"        errmsg:"Please provide the invoice id"`
	PaymentIntent       string            `json:"paymentIntent"               bson:"paymentIntent"         validate:"required"        errmsg:"Please provide payment intent"`
	GigId               string            `json:"gigId"                       bson:"gigId"                 validate:"required"        errmsg:"Please provide gig id"`
	SellerId            string            `json:"sellerId"                    bson:"sellerId"              validate:"required"        errmsg:"Please provide seller id"`
	SellerUsername      string            `json:"sellerUsername"              bson:"sellerUsername"        validate:"required"        errmsg:"Please provide seller username"`
	SellerImage         string            `json:"sellerImage"                 bson:"sellerImage"           validate:"required,url"    errmsg:"Please provide seller image"`
	SellerEmail         string            `json:"sellerEmail"                 bson:"sellerEmail"           validate:"required,email"  errmsg:"Please provide seller email"`
	GigCoverImage       string            `json:"gigCoverImage"               bson:"gigCoverImage"         validate:"required,url"    errmsg:"Please provide gig image"`
	GigMainTitle        string            `json:"gigMainTitle"                bson:"gigMainTitle"          validate:"required"        errmsg:"Please provide gig main title"`
	GigBasicTitle       string            `json:"gigBasicTitle"               bson:"gigBasicTitle"         validate:"required"        errmsg:"Please provide gig basic title"`
	GigBasicDescription string            `json:"gigBasicDescription"         bson:"gigBasicDescription"   validate:"required"        errmsg:"Please provide gig basic description"`
	BuyerId             string            `json:"buyerId"                     bson:"buyerId"               validate:"required"        errmsg:"Please provide buyer id"`
	BuyerUsername       string            `json:"buyerUsername"               bson:"buyerUsername"         validate:"required"        errmsg:"Please provide buyer username"`
	BuyerEmail          string            `json:"buyerEmail"                  bson:"buyerEmail"            validate:"required,email"  errmsg:"Please provide buyer email"`
	BuyerImage          string            `json:"buyerImage"                  bson:"buyerImage"            validate:"required,url"    errmsg:"Please provide buyer image"`
	Status              string            `json:"status"                      bson:"status"                validate:"required"        errmsg:"Please provide status"`
	Requirements        string            `json:"requirements"                bson:"requirements"          `
	Quantity            int32             `json:"quantity"                    bson:"quantity"              validate:"required,gt=0"   errmsg:"Please provide quantity"`
	Price               float32           `json:"price"                       bson:"price"                 validate:"required,gt=0"   errmsg:"Please provide price"`
	ServiceFee          float32           `json:"serviceFee"                  bson:"serviceFee"            `
	Approved            bool              `json:"approved"                    bson:"approved"              `
	Delivered           bool              `json:"delivered"                   bson:"delivered"             `
	Cancelled           bool              `json:"cancelled"                   bson:"cancelled"             `
	ApprovedAt          *time.Time        `json:"approvedAt"                  bson:"approvedAt"            `
	DateOrdered         *time.Time        `json:"dateOrdered"                 bson:"dateOrdered"           `
	Offer               *Offer            `json:"offer"                       bson:"offer"                 validate:"required"        errmsg:"Please provide offer data"`
	DeliveredWork       []*DeliveredWork  `json:"deliveredWork"               bson:"deliveredWork"         `
	Events              *OrderEvents      `json:"events"                      bson:"events"                `
	BuyerReview         *OrderReview      `json:"buyerReview,omitempty"       bson:"buyerReview"           `
	SellerReview        *OrderReview      `json:"sellerReview,omitempty"      bson:"sellerReview"          `
	RequestExtension    *ExtendedDelivery `json:"requestExtension,omitempty"  bson:"requestExtension"      `
}

type DeliveredWork struct {
	Message  string `json:"message"    bson:"message"   validate:"required"  errmsg:"Please provide the message for the delivery"`
	File     string `json:"file"       bson:"file"      validate:"required"  errmsg:"Please provide the file data"`
	FileType string `json:"fileType"   bson:"fileType"  validate:"required"  errmsg:"Please provide the file type"`
	FileSize int64  `json:"fileSize"   bson:"fileSize"  validate:"required"  errmsg:"Please provide the file size"`
	FileName string `json:"fileName"   bson:"fileName"  validate:"required"  errmsg:"Please provide the file name"`
}

type OrderEvents struct {
	PlaceOrder         string `json:"placeOrder"          bson:"placeOrder"`
	Requirements       string `json:"requirements"        bson:"requirements"`
	OrderStarted       string `json:"orderStarted"        bson:"orderStarted"`
	DeliveryDateUpdate string `json:"deliveryDateUpdate"  bson:"deliveryDateUpdate"`
	OrderDelivered     string `json:"orderDelivered"      bson:"orderDelivered"`
	BuyerReview        string `json:"buyerReview"         bson:"buyerReview"`
	SellerReview       string `json:"sellerReview"        bson:"sellerReview"`
}

type OrderReview struct {
	Rating int32  `json:"rating"  bson:"rating"`
	Review string `json:"review"  bson:"review"`
	Date   string `json:"date"    bson:"date"`
}

type ExtendedDelivery struct {
	OriginalDate       *time.Time `json:"originalDate"         bson:"originalDate"         validate:"required"  errmsg:"Please provide the original date"`
	NewDate            *time.Time `json:"newDate"              bson:"newDate"              validate:"required"  errmsg:"Please provide the new date"`
	Days               int32      `json:"days"                 bson:"days"                 validate:"required"  errmsg:"Please provide the number of days"`
	Reason             string     `json:"reason"               bson:"reason"               validate:"required"  errmsg:"Please provide a reason"`
	DeliveryDateUpdate string     `json:"deliveryDateUpdate"   bson:"deliveryDateUpdate"   `
}

type OrderResponse struct {
	Message string         `json:"message"`
	Order   *OrderDocument `json:"order"`
}

type PaymentIntentResponse struct {
	Message         string `json:"message"`
	ClientSecret    string `json:"clientSecret"`
	PaymentIntentId string `json:"paymentIntentId"`
}

type OrdersResponse struct {
	Message string           `json:"message"`
	Orders  []*OrderDocument `json:"orders"`
}

type OrderMessage struct {
	SellerId       string  `json:"sellerId"`
	BuyerId        string  `json:"buyerId"`
	OngoingJobs    int32   `json:"ongoingJobs"`
	CompletedJobs  int32   `json:"completedJobs"`
	TotalEarnings  float32 `json:"totalEarnings"`
	PurchasedGigs  string  `json:"purchasedGigs"`
	RecentDelivery string  `json:"recentDelivery"`
	Type           string  `json:"type"`
	ReceiverEmail  string  `json:"receiverEmail"`
	Username       string  `json:"username"`
	Template       string  `json:"template"`
	Sender         string  `json:"sender"`
	OfferLink      string  `json:"offerLink"`
	Amount         string  `json:"amount"`
	BuyerUsername  string  `json:"buyerUsername"`
	SellerUsername string  `json:"sellerUsername"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	DeliveryDays   string  `json:"deliveryDays"`
	OrderId        string  `json:"orderId"`
	InvoiceId      string  `json:"invoiceId"`
	OrderDue       string  `json:"orderDue"`
	Requirements   string  `json:"requirements"`
	OrderUrl       string  `json:"orderUrl"`
	OriginalDate   string  `json:"originalDate"`
	NewDate        string  `json:"newDate"`
	Reason         string  `json:"reason"`
	Subject        string  `json:"subject"`
	Header         string  `json:"header"`
	Total          string  `json:"total"`
	Message        string  `json:"message"`
	ServiceFee     string  `json:"serviceFee"`
}

type CancelOrderRequest struct {
	PaymentIntent string        `json:"paymentIntent"`
	OrderData     *OrderMessage `json:"orderData"`
}
