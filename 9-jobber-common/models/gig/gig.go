package gig

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/thetherington/jobber-common/models/review"
	pb "github.com/thetherington/jobber-common/protogen/go/gig"
	pbReview "github.com/thetherington/jobber-common/protogen/go/review"
	"github.com/thetherington/jobber-common/utils"
)

type SellerGig struct {
	ES_Id            string                   `json:"-"                           bson:"_id,omitempty"     `
	ID               string                   `json:"id,omitempty"                bson:"-"                 `
	SellerId         string                   `json:"sellerId,omitempty"          bson:"sellerId"          validate:"required"        errmsg:"Please provide the seller id"`
	Title            string                   `json:"title"                       bson:"title"             validate:"required"        errmsg:"Please provide a title"`
	Username         string                   `json:"username,omitempty"          bson:"username"          `
	ProfilePicture   string                   `json:"profilePicture,omitempty"    bson:"profilePicture"    validate:"required,url"    errmsg:"Please provide a profile picture"`
	Email            string                   `json:"email,omitempty"             bson:"email"             `
	Description      string                   `json:"description"                 bson:"description"       validate:"required"        errmsg:"Please provide a description"`
	Active           bool                     `json:"active"                      bson:"active"            `
	Categories       string                   `json:"categories"                  bson:"categories"        validate:"required"        errmsg:"Please provide the category"`
	SubCategories    []string                 `json:"subCategories"               bson:"subCategories"     validate:"required,min=1"  errmsg:"Please provide atleast 1 sub category"`
	Tags             []string                 `json:"tags"                        bson:"tags"              validate:"required,min=1"  errmsg:"Please provide atleast 1 tag"`
	RatingsCount     int32                    `json:"ratingsCount"                bson:"ratingsCount"      `
	RatingSum        int32                    `json:"ratingSum"                   bson:"ratingSum"         `
	RatingCategories *review.RatingCategories `json:"ratingCategories,omitempty"  bson:"ratingCategories"  `
	ExpectedDelivery string                   `json:"expectedDelivery"            bson:"expectedDelivery"  validate:"required,min=1"  errmsg:"Please provide atleast 1 tag"`
	BasicTitle       string                   `json:"basicTitle"                  bson:"basicTitle"        validate:"required"        errmsg:"Please provide a basic title"`
	BasicDescription string                   `json:"basicDescription"            bson:"basicDescription"  validate:"required"        errmsg:"Please provide a basic description"`
	Price            float32                  `json:"price"                       bson:"price"             validate:"required,min=1"  errmsg:"Please provide a price greater than $1"`
	CoverImage       string                   `json:"coverImage"                  bson:"coverImage"        validate:"required"        errmsg:"Please provide a cover image"`
	CreatedAt        *time.Time               `json:"createdAt,omitempty"         bson:"createdAt"         `
	SortId           int32                    `json:"sortId"                      bson:"sortId"            `
}

type UpdateSellerGig struct {
	Title            string   `json:"title"                       bson:"title"             validate:"required"        errmsg:"Please provide a title"`
	Description      string   `json:"description"                 bson:"description"       validate:"required"        errmsg:"Please provide a description"`
	Categories       string   `json:"categories"                  bson:"categories"        validate:"required"        errmsg:"Please provide the category"`
	SubCategories    []string `json:"subCategories"               bson:"subCategories"     validate:"required,min=1"  errmsg:"Please provide atleast 1 sub category"`
	Tags             []string `json:"tags"                        bson:"tags"              validate:"required,min=1"  errmsg:"Please provide atleast 1 tag"`
	ExpectedDelivery string   `json:"expectedDelivery"            bson:"expectedDelivery"  validate:"required,min=1"  errmsg:"Please provide atleast 1 tag"`
	BasicTitle       string   `json:"basicTitle"                  bson:"basicTitle"        validate:"required"        errmsg:"Please provide a basic title"`
	BasicDescription string   `json:"basicDescription"            bson:"basicDescription"  validate:"required"        errmsg:"Please provide a basic description"`
	Price            float32  `json:"price"                       bson:"price"             validate:"required,min=1"  errmsg:"Please provide a price greater than $1"`
}

func (s *SellerGig) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[SellerGig](*s, validate)
}

func (s *UpdateSellerGig) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[UpdateSellerGig](*s, validate)
}

type ResponseSearchGigs struct {
	Message string       `json:"message"`
	Total   int32        `json:"total"`
	Gigs    []*SellerGig `json:"gigs"`
}

type ResponseGig struct {
	Message string     `json:"message"`
	Gig     *SellerGig `json:"gig"`
}

type ResponseGigs struct {
	Message string       `json:"message"`
	Gigs    []*SellerGig `json:"gigs"`
}

func CreateSellerGig(gigpb *pb.GigMessage) *SellerGig {
	g := SellerGig{
		ES_Id:            *gigpb.ES_ID,
		ID:               *gigpb.ID,
		SellerId:         *gigpb.SellerId,
		Title:            *gigpb.Title,
		Username:         *gigpb.Username,
		ProfilePicture:   *gigpb.ProfilePicture,
		Email:            *gigpb.Email,
		Description:      gigpb.Description,
		Active:           gigpb.Active,
		Categories:       gigpb.Categories,
		SubCategories:    gigpb.SubCategories,
		Tags:             gigpb.Tags,
		RatingsCount:     gigpb.RatingsCount,
		RatingSum:        gigpb.RatingSum,
		ExpectedDelivery: gigpb.ExpectedDelivery,
		BasicTitle:       gigpb.BasicTitle,
		BasicDescription: gigpb.BasicDescription,
		Price:            gigpb.Price,
		CoverImage:       gigpb.CoverImage,
		CreatedAt:        utils.ToTime(gigpb.GetCreatedAt()),
		SortId:           gigpb.SortId,
	}

	if gigpb.RatingCategories != nil {
		rc := gigpb.RatingCategories

		g.RatingCategories = &review.RatingCategories{
			One:   review.RatingCategoryItem{Value: rc.One.Value, Count: rc.One.Count},
			Two:   review.RatingCategoryItem{Value: rc.Two.Value, Count: rc.Two.Count},
			Three: review.RatingCategoryItem{Value: rc.Three.Value, Count: rc.Three.Count},
			Four:  review.RatingCategoryItem{Value: rc.Four.Value, Count: rc.Four.Count},
			Five:  review.RatingCategoryItem{Value: rc.Five.Value, Count: rc.Five.Count},
		}
	}

	return &g
}

func CreateGigMessage(gig *SellerGig) *pb.GigMessage {
	g := pb.GigMessage{
		ES_ID:            &gig.ES_Id,
		ID:               &gig.ID,
		SellerId:         &gig.SellerId,
		Title:            &gig.Title,
		Username:         &gig.Username,
		ProfilePicture:   &gig.ProfilePicture,
		Email:            &gig.Email,
		Description:      gig.Description,
		Active:           gig.Active,
		Categories:       gig.Categories,
		SubCategories:    gig.SubCategories,
		Tags:             gig.Tags,
		RatingsCount:     gig.RatingsCount,
		RatingSum:        gig.RatingSum,
		ExpectedDelivery: gig.ExpectedDelivery,
		BasicTitle:       gig.BasicTitle,
		BasicDescription: gig.BasicDescription,
		Price:            gig.Price,
		CoverImage:       gig.CoverImage,
		CreatedAt:        utils.ToDateTime(gig.CreatedAt),
		SortId:           gig.SortId,
	}

	if gig.RatingCategories != nil {
		rc := gig.RatingCategories

		g.RatingCategories = &pbReview.RatingCategories{
			One:   &pbReview.RatingCategoryItem{Value: rc.One.Value, Count: rc.One.Count},
			Two:   &pbReview.RatingCategoryItem{Value: rc.Two.Value, Count: rc.Two.Count},
			Three: &pbReview.RatingCategoryItem{Value: rc.Three.Value, Count: rc.Three.Count},
			Four:  &pbReview.RatingCategoryItem{Value: rc.Four.Value, Count: rc.Four.Count},
			Five:  &pbReview.RatingCategoryItem{Value: rc.Five.Value, Count: rc.Five.Count},
		}
	}

	return &g
}
