package notification

type AuthEmailLocals struct {
	AppLink       string
	AppIcon       string
	Username      string
	VerifyLink    string
	ResetLink     string
	ReceiverEmail string
}

type OrderEmailLocals struct {
	AppLink        string
	AppIcon        string
	ReceiverEmail  string
	Username       string
	Template       string
	Sender         string
	OfferLink      string
	Amount         string
	BuyerUsername  string
	SellerUsername string
	Title          string
	Description    string
	DeliveryDays   string
	InvoiceId      string
	OrderId        string
	OrderDue       string
	Requirements   string
	OrderUrl       string
	OriginalDate   string
	NewDate        string
	Reason         string
	Subject        string
	Header         string
	Type           string
	Message        string
	ServiceFee     string
	Total          string
}
