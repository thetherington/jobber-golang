package port

import (
	"context"

	"github.com/thetherington/jobber-common/models/review"
	"github.com/thetherington/jobber-common/models/users"
)

type BuyerService interface {
	CreateBuyer(b *users.Buyer) error
	GetBuyerByEmail(ctx context.Context, email string) (*users.Buyer, error)
	GetBuyerByUsername(ctx context.Context, username string) (*users.Buyer, error)
	GetRandomBuyers(ctx context.Context, count int) ([]*users.Buyer, error)
	UpdateBuyerIsSeller(ctx context.Context, email string) error
	UpdateBuyerPurchasedGigs(buyerId string, pruchasedGigId string, action string) error
}

type SellerService interface {
	CreateSeller(ctx context.Context, seller *users.Seller) (*users.SellerResponse, error)
	UpdateSeller(ctx context.Context, id string, req *users.Seller) (*users.SellerResponse, error)
	GetSellerById(ctx context.Context, id string) (*users.SellerResponse, error)
	GetSellerByUsername(ctx context.Context, username string) (*users.SellerResponse, error)
	GetRandomSellers(ctx context.Context, count int32) (*users.SellersResponse, error)
	SeedSellers(ctx context.Context, count int32) (string, error)
	UpdateTotalGigCount(id string, count int32) error
	UpdateSellerOngoingJobsProp(id string, ongoingjobs int32) error
	UpdateSellerCancelledJobsProp(id string) error
	UpdateSellerCompletedJobsProp(sellerId string, ongoingJobs int32, completedJobs int32, totalEarnings float32) error
	UpdateSellerReview(data *review.ReviewMessageDetails) error
}
