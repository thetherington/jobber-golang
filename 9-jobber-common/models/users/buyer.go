package users

import "time"

type Buyer struct {
	Id             string     `json:"_id"             bson:"_id,omitempty"`
	Username       string     `json:"username"        bson:"username"`
	Email          string     `json:"email"           bson:"email"`
	ProfilePicture string     `json:"profilePicture"  bson:"profilePicture"`
	Country        string     `json:"country"         bson:"country"`
	IsSeller       bool       `json:"isSeller"        bson:"isSeller"`
	PurchasedGigs  []string   `json:"purchasedGigs"   bson:"purchasedGigs"`
	CreatedAt      *time.Time `json:"createdAt"       bson:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"       bson:"updatedAt"`
}

type BuyerResponse struct {
	Message string `json:"message"`
	Buyer   *Buyer `json:"buyer"`
}
