package review

import (
	"time"

	"github.com/go-playground/validator/v10"
	pb "github.com/thetherington/jobber-common/protogen/go/review"
	"github.com/thetherington/jobber-common/utils"
)

type RatingCategoryItem struct {
	Value int32 `json:"value"`
	Count int32 `json:"count"`
}

type RatingCategories struct {
	Five  RatingCategoryItem `json:"five"`
	Four  RatingCategoryItem `json:"four"`
	Three RatingCategoryItem `json:"three"`
	Two   RatingCategoryItem `json:"two"`
	One   RatingCategoryItem `json:"one"`
}

type ReviewMessageDetails struct {
	GigId      string    `json:"gigId"`
	ReviewerId string    `json:"reviewerId"`
	SellerId   string    `json:"sellerId"`
	Review     string    `json:"review"`
	Rating     int32     `json:"rating"`
	OrderId    string    `json:"orderId"`
	CreatedAt  time.Time `json:"createdAt"`
	Type       string    `json:"type"`
}

type ReviewDocument struct {
	Id               string    `json:"id"`
	GigId            string    `json:"gigId"            validate:"required"        errmsg:"Please provide the gig id"`
	ReviewerId       string    `json:"reviewerId"       validate:"required"        errmsg:"Please provide the reviewer id"`
	OrderId          string    `json:"orderId"          validate:"required"        errmsg:"Please provide the order id"`
	SellerId         string    `json:"sellerId"         validate:"required"        errmsg:"Please provide the seller id"`
	Review           string    `json:"review"           validate:"required"        errmsg:"Please provide a review"`
	ReviewerImage    string    `json:"reviewerImage"    validate:"required,url"    errmsg:"Please provide the reviewer url image"`
	ReviewerUsername string    `json:"reviewerUsername" validate:"required"        errmsg:"Please provide username of the reviewer"`
	Country          string    `json:"country"          validate:"required"        errmsg:"Please provide a country"`
	ReviewType       string    `json:"reviewType"       validate:"required"        errmsg:"Please provide a review type"`
	Rating           int32     `json:"rating"           validate:"required,gte=1"  errmsg:"Please provide a rating"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (s *ReviewDocument) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[ReviewDocument](*s, validate)
}

type ReviewResponse struct {
	Message string          `json:"message"`
	Review  *ReviewDocument `json:"review"`
}

type ReviewsResponse struct {
	Message string            `json:"message"`
	Reviews []*ReviewDocument `json:"reviews"`
}

func (r *ReviewDocument) MarshalToProto() *pb.ReviewDocument {
	return &pb.ReviewDocument{
		Id:               r.Id,
		GigId:            r.GigId,
		ReviewerId:       r.ReviewerId,
		OrderId:          r.OrderId,
		SellerId:         r.SellerId,
		Review:           r.Review,
		ReviewerImage:    r.ReviewerImage,
		ReviewerUsername: r.ReviewerUsername,
		Country:          r.Country,
		ReviewType:       r.ReviewType,
		Rating:           r.Rating,
		CreatedAt:        utils.ToDateTime(&r.CreatedAt),
	}
}

func UnmarshalToDocument(r *pb.ReviewDocument) *ReviewDocument {
	return &ReviewDocument{
		Id:               r.Id,
		GigId:            r.GigId,
		ReviewerId:       r.ReviewerId,
		OrderId:          r.OrderId,
		SellerId:         r.SellerId,
		Review:           r.Review,
		ReviewerImage:    r.ReviewerImage,
		ReviewerUsername: r.ReviewerUsername,
		Country:          r.Country,
		ReviewType:       r.ReviewType,
		Rating:           r.Rating,
		CreatedAt:        *utils.ToTime(r.CreatedAt),
	}
}
