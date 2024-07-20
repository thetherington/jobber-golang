package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/middleware"
	"github.com/thetherington/jobber-common/models/review"
	"github.com/thetherington/jobber-common/models/users"
	"github.com/thetherington/jobber-common/utils"
	"github.com/thetherington/jobber-users/internal/core/port"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	sellerDB *mongo.Collection
	validate *validator.Validate
)

/**
 * SellerService implements
 */
type SellerService struct {
	buyer port.BuyerService
}

// NewSellerService creates a new buyer service instance
func NewSellerService(db *mongo.Database, buyer port.BuyerService) *SellerService {
	validate = validator.New(validator.WithRequiredStructEnabled())

	sellerDB = db.Collection("Seller")

	sellerDB.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}},
		{Keys: bson.D{{Key: "email", Value: 1}}},
	})

	return &SellerService{
		buyer: buyer,
	}
}

func (s *SellerService) CreateSeller(ctx context.Context, seller *users.Seller) (*users.SellerResponse, error) {
	// check that the seller with the email doesn't already exists.
	if s, _ := s.GetSellerByEmail(ctx, seller.Email); s != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("seller already exists with the email %s", seller.Email))
	}

	// Validate seller payload
	if err := seller.Validate(validate); err != nil {
		slog.With("error", err).Debug("Seller Create Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// Get username from the user cookie session passed down into the context.
	username := ctx.Value(middleware.CtxUsernameKey)
	if username == nil {
		slog.Debug("Username in context is nil")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again"))
	}

	// only set the username from the cookie session if the username is blank
	if seller.Username == "" {
		seller.Username = username.(string)
	}

	// insert seller into Sellers mongodb collection
	resp, err := sellerDB.InsertOne(ctx, seller, &options.InsertOneOptions{})
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// update buyer to be seller
	if err := s.buyer.UpdateBuyerIsSeller(ctx, seller.Email); err != nil {
		slog.With("error", err).Error("failed to update buyer to be seller by email", "email", seller.Email)
	}

	objId, ok := resp.InsertedID.(primitive.ObjectID)
	if ok {
		seller.Id = objId.Hex()
	}

	return &users.SellerResponse{Message: "Seller created successfully", Seller: seller}, nil
}

func (s *SellerService) UpdateSeller(ctx context.Context, id string, seller *users.Seller) (*users.SellerResponse, error) {
	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	resp, err := sellerDB.UpdateByID(ctx, objectId,
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "profilePublicId", Value: seller.ProfilePublicId},
			{Key: "fullName", Value: seller.FullName},
			{Key: "profilePicture", Value: seller.ProfilePicture},
			{Key: "description", Value: seller.Description},
			{Key: "country", Value: seller.Country},
			{Key: "skills", Value: seller.Skills},
			{Key: "oneliner", Value: seller.Oneliner},
			{Key: "languages", Value: seller.Languages},
			{Key: "responseTime", Value: seller.ResponseTime},
			{Key: "experience", Value: seller.Experience},
			{Key: "education", Value: seller.Education},
			{Key: "socialLinks", Value: seller.SocialLinks},
			{Key: "certificates", Value: seller.Certificates},
		}}},
	)
	if err != nil {
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}
	if resp.MatchedCount == 0 {
		return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("seller id does not exist"))
	}
	if resp.ModifiedCount == 0 {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("seller was not updated"))
	}

	return &users.SellerResponse{Message: "Seller updated successfully", Seller: seller}, nil
}

func (s *SellerService) GetSellerById(ctx context.Context, id string) (*users.SellerResponse, error) {
	var seller *users.Seller

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	err = sellerDB.FindOne(ctx, bson.M{"_id": objectId}).Decode(&seller)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("seller id does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &users.SellerResponse{Message: "Seller", Seller: seller}, nil
}

func (s *SellerService) GetSellerByUsername(ctx context.Context, username string) (*users.SellerResponse, error) {
	var seller *users.Seller

	err := sellerDB.FindOne(ctx, bson.D{{Key: "username", Value: utils.FirstLetterUpperCase(username)}}).Decode(&seller)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("seller username does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}
	return &users.SellerResponse{Message: "Seller", Seller: seller}, nil
}

func (s *SellerService) GetSellerByEmail(ctx context.Context, email string) (*users.SellerResponse, error) {
	var seller *users.Seller

	err := sellerDB.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&seller)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, svc.NewError(svc.ErrNotFound, fmt.Errorf("seller email does not exist"))
		}

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	return &users.SellerResponse{Message: "Seller", Seller: seller}, nil
}

func (s *SellerService) GetRandomSellers(ctx context.Context, count int32) (*users.SellersResponse, error) {
	agg := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: count}}}}

	cursor, err := sellerDB.Aggregate(ctx, mongo.Pipeline{agg})
	if err != nil {
		return nil, err
	}

	var sellers []*users.Seller

	if err = cursor.All(context.TODO(), &sellers); err != nil {
		return nil, err
	}

	return &users.SellersResponse{
		Message: "Random sellers profile",
		Sellers: sellers,
	}, nil
}

func (s *SellerService) SeedSellers(ctx context.Context, count int32) (string, error) {
	buyers, err := s.buyer.GetRandomBuyers(ctx, int(count))
	if err != nil {
		return "", svc.NewError(svc.ErrInternalFailure, err)
	}

	var (
		f             = gofakeit.New(0)
		skills        = []string{"Programming", "Web development", "Mobile development", "Proof reading", "UI/UX", "Data Science", "Financial modeling", "Data analysis"}
		LanguageLevel = []string{"Native", "Basic", "Advance"}
	)

	randomSkills := func(count int) []string {
		var x []string
		for i := 0; i < count; i++ {
			x = append(x, f.RandomString(skills))
		}
		return x
	}

	randomExperience := func(count int) []users.Experience {
		var (
			x               = make([]users.Experience, 0)
			randomStartYear = []string{"2020", "2021", "2022", "2023", "2024", "2025"}
			randomEndYear   = []string{"Present", "2025", "2026", "2027"}
			endYear         = f.RandomString(randomEndYear)
		)
		endYearString := func(y string) string {
			if y == "Present" {
				return y
			}
			return fmt.Sprintf("%s %s", f.MonthString(), f.RandomString(randomEndYear))
		}
		endYearBool := func(y string) bool {
			return y == "Present"
		}
		for i := 0; i < count; i++ {
			e := users.Experience{
				Id:                   primitive.NewObjectID().Hex(),
				Company:              f.Company(),
				Title:                f.JobTitle(),
				Description:          f.JobDescriptor(),
				StartDate:            fmt.Sprintf("%s %s", f.MonthString(), f.RandomString(randomStartYear)),
				EndDate:              endYearString(endYear),
				CurrentlyWorkingHere: endYearBool(endYear),
			}
			x = append(x, e)
		}
		return x
	}

	randomEducation := func(count int) []users.Education {
		var (
			x          []users.Education
			randomYear = []string{"2020", "2021", "2022", "2023", "2024", "2025"}
		)
		for i := 0; i < count; i++ {
			e := users.Education{
				Id:         primitive.NewObjectID().Hex(),
				Country:    f.Country(),
				University: f.School(),
				Title:      f.JobTitle(),
				Major:      f.JobDescriptor(),
				Year:       f.RandomString(randomYear),
			}
			x = append(x, e)
		}
		return x
	}

	for _, buyer := range buyers {
		// check that the seller with the email doesn't already exists.
		if s, _ := s.GetSellerByEmail(ctx, buyer.Email); s != nil {
			return "", svc.NewError(svc.ErrBadRequest, fmt.Errorf("seller already exists with the email %s", buyer.Email))
		}

		recentDelivery := f.Date()

		seller := &users.Seller{
			ProfilePublicId: uuid.New().String(),
			FullName:        fmt.Sprintf("%s %s", f.FirstName(), f.LastName()),
			Username:        buyer.Username,
			Email:           buyer.Email,
			ProfilePicture:  buyer.ProfilePicture,
			Description:     f.Sentence(25),
			Country:         f.Country(),
			Oneliner:        f.Sentence(10),
			Skills:          randomSkills(f.Number(2, 4)),
			ResponseTime:    int32(f.Number(1, 5)),
			Languages: []users.Language{
				{Id: primitive.NewObjectID().Hex(), Language: f.Language(), Level: f.RandomString(LanguageLevel)},
				{Id: primitive.NewObjectID().Hex(), Language: f.Language(), Level: f.RandomString(LanguageLevel)},
				{Id: primitive.NewObjectID().Hex(), Language: f.Language(), Level: f.RandomString(LanguageLevel)},
			},
			RatingCategories: review.RatingCategories{
				One:   review.RatingCategoryItem{},
				Two:   review.RatingCategoryItem{},
				Three: review.RatingCategoryItem{},
				Four:  review.RatingCategoryItem{},
				Five:  review.RatingCategoryItem{},
			},
			SocialLinks:    []string{"https://kickchatapp.com", "http://youtube.com", "https://facebook.com"},
			RecentDelivery: &recentDelivery,
			Experience:     randomExperience(f.Number(2, 4)),
			Education:      randomEducation(f.Number(2, 4)),
			Certificates: []users.Certificate{
				{Id: primitive.NewObjectID().Hex(), Name: "Flutter App Develope", From: "Flutter Academy", Year: "2021"},
				{Id: primitive.NewObjectID().Hex(), Name: "Android App Developer", From: "2019", Year: "2020"},
				{Id: primitive.NewObjectID().Hex(), Name: "IOS App Developer", From: "Apple Inc.", Year: "2019"},
			},
			RatingsCount:  0,
			RatingSum:     0,
			OngoingJobs:   0,
			CompletedJobs: 0,
			CancelledJobs: 0,
			TotalEarnings: 0,
			TotalGigs:     0,
			CreatedAt:     f.Date(),
			UpdatedAt:     f.Date(),
		}

		_, err := s.CreateSeller(ctx, seller)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("Seeded %d sellers successfully", count), nil
}

func (s *SellerService) UpdateTotalGigCount(id string, count int32) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := sellerDB.UpdateByID(ctx, objectId, bson.D{{Key: "$inc", Value: bson.D{{Key: "totalGigs", Value: count}}}})
	return UpdateErrorCheck(resp, err)
}

func (s *SellerService) UpdateSellerOngoingJobsProp(id string, ongoingjobs int32) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := sellerDB.UpdateByID(ctx, objectId, bson.D{{Key: "$inc", Value: bson.D{{Key: "ongoingJobs", Value: ongoingjobs}}}})
	return UpdateErrorCheck(resp, err)
}

func (s *SellerService) UpdateSellerCancelledJobsProp(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := sellerDB.UpdateByID(ctx, objectId, bson.D{
		{Key: "$inc",
			Value: bson.D{
				{Key: "ongoingJobs", Value: -1},
				{Key: "cancelledJobs", Value: 1},
			},
		},
	})
	return UpdateErrorCheck(resp, err)
}

func (s *SellerService) UpdateSellerCompletedJobsProp(sellerId string, ongoingJobs int32, completedJobs int32, totalEarnings float32) error {
	objectId, err := primitive.ObjectIDFromHex(sellerId) //sellerId
	if err != nil {
		return svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := sellerDB.UpdateByID(ctx, objectId, bson.D{
		{Key: "$inc",
			Value: bson.D{
				{Key: "ongoingJobs", Value: ongoingJobs},     //ongoingjobs
				{Key: "completedJobs", Value: completedJobs}, //completedjobs
				{Key: "totalEarnings", Value: totalEarnings}, //totalEarnings
			},
		},
		{Key: "$set",
			Value: bson.D{
				{Key: "recentDelivery", Value: time.Now()}, //recentDelivery
			},
		},
	})
	return UpdateErrorCheck(resp, err)
}

func (s *SellerService) UpdateSellerReview(data *review.ReviewMessageDetails) error {
	objectId, err := primitive.ObjectIDFromHex(data.SellerId)
	if err != nil {
		return svc.NewError(svc.ErrBadRequest, fmt.Errorf("id provided is invalid"))
	}

	ratingTypes := map[string]string{
		"1": "one",
		"2": "two",
		"3": "three",
		"4": "four",
		"5": "five",
	}
	ratingKey := ratingTypes[strconv.Itoa(int(data.Rating))]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := sellerDB.UpdateByID(ctx, objectId, bson.D{
		{Key: "$inc",
			Value: bson.D{
				{Key: "ratingsCount", Value: 1},
				{Key: "ratingSum", Value: data.Rating},
				{Key: fmt.Sprintf("ratingCategories.%s.value", ratingKey), Value: data.Rating},
				{Key: fmt.Sprintf("ratingCategories.%s.count", ratingKey), Value: 1},
			},
		},
	})

	return UpdateErrorCheck(resp, err)
}
