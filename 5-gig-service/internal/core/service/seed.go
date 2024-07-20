package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/thetherington/jobber-common/models/gig"
	"github.com/thetherington/jobber-common/models/review"
	pborder "github.com/thetherington/jobber-common/protogen/go/order"
	pb "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
)

func (g *GigService) SeedGigs(ctx context.Context, count int32) (string, error) {
	// send the user information to the notification micro service via rabbitMQ direct exchange using protobuf
	if data, err := proto.Marshal(&pb.RandomSellersRequest{Size: count}); err == nil {
		if err := g.queue.PublishDirectMessage("jobber-gig", "get-sellers", data); err != nil {
			slog.With("error", err).Error("SeedGig: Failed to send message to jobber-gig")
		}
	}

	return "Gigs created successfully", nil
}

func (g *GigService) SeedData(ctx context.Context, sellers []any) error {
	f := gofakeit.New(0)

	categories := []string{
		"Graphics & Design",
		"Digital Marketing",
		"Writing & Translation",
		"Video & Animation",
		"Music & Audio",
		"Programming & Tech",
		"Data",
		"Business",
	}

	expectedDelivery := []string{"1 Day Delivery", "2 Days Delivery", "3 Days Delivery", "4 Days Delivery", "5 Days Delivery"}

	type rating struct {
		Sum   int
		Count int
	}
	// get gigs count from elastic
	count, err := g.search.GetGigsCount(ctx)
	if err != nil {
		slog.With("error", err).Error("failed to get count")
	}

	randomRatings := []rating{
		{Sum: 20, Count: 4},
		{Sum: 10, Count: 2},
		{Sum: 20, Count: 4},
		{Sum: 15, Count: 3},
		{Sum: 5, Count: 1},
	}

	for i, s := range sellers {
		pbSeller, ok := s.(*pb.SellerPayload)
		if !ok {
			continue
		}

		rating := randomRatings[f.IntN(len(randomRatings))]

		gig := &gig.SellerGig{
			SellerId:         pbSeller.Id,
			ProfilePicture:   pbSeller.ProfilePicture,
			Email:            pbSeller.Email,
			Username:         pbSeller.Username,
			Title:            fmt.Sprintf("I will %s", f.Sentence(5)),
			BasicTitle:       f.ProductName(),
			BasicDescription: f.ProductDescription(),
			Description:      f.LoremIpsumSentence(25),
			Categories:       f.RandomString(categories),
			SubCategories:    []string{f.ProductCategory(), f.ProductCategory(), f.ProductCategory()},
			Tags:             []string{f.ProgrammingLanguage(), f.ProgrammingLanguage(), f.ProgrammingLanguage()},
			Price:            float32(f.Price(20, 30)),
			ExpectedDelivery: f.RandomString(expectedDelivery),
			CoverImage:       fmt.Sprintf("https://picsum.photos/seed/%s/640/480", randomString(7)),
			SortId:           int32(int(count) + i + 1),
			RatingSum:        int32(rating.Sum),
			RatingsCount:     int32(rating.Count),
			Active:           true,
			CreatedAt:        utils.ToTimePtr(f.Date()),
			RatingCategories: &review.RatingCategories{
				One:   review.RatingCategoryItem{Value: 0, Count: 0},
				Two:   review.RatingCategoryItem{Value: 0, Count: 0},
				Three: review.RatingCategoryItem{Value: 0, Count: 0},
				Four:  review.RatingCategoryItem{Value: 0, Count: 0},
				Five:  review.RatingCategoryItem{Value: 0, Count: 0},
			},
		}

		g.SeedGigCreate(ctx, gig)
	}

	return nil
}

func (g *GigService) SeedGigCreate(ctx context.Context, gig *gig.SellerGig) {
	// save gig into mongodb, get _id and use it in id
	result, err := gigDB.InsertOne(ctx, gig)
	if err != nil {
		slog.With("error", err).Error("failed to insert gig into mongodb")
		return
	}

	if objId, ok := result.InsertedID.(primitive.ObjectID); ok {
		gig.ID = objId.Hex()
	}

	// save gig into elasticsearch
	if _, err = g.search.InsertGig(ctx, gig); err != nil {
		slog.With("error", err).Error("failed to insert gig elastic")
		return
	}

	// increment the seller gig count via the Users microservice from rabbitmq
	pbmsg := &pborder.SellerGigUpdate{
		Action:     pborder.Action_UpdateGigCount,
		SellerId:   gig.SellerId,
		OrderProps: &pborder.OrderProps{GigCount: proto.Int32(1)},
	}

	if data, err := proto.Marshal(pbmsg); err == nil {
		if err := g.queue.PublishDirectMessage("jobber-seller-update", "user-seller", data); err != nil {
			slog.With("error", err).Error("Signup: Failed to send message to jobber-seller-update")
		}
	}

}

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString returns a string of random characters of length n, using randomStringSource
// as a source for the string
func randomString(n int) string {

	s, r := make([]rune, n), []rune(randomStringSource)

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}
