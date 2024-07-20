package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/thetherington/jobber-common/models/review"
	"github.com/thetherington/jobber-review/internal/adapters/storage/postgres"
)

/**
 * ReviewRepository implements port.ReviewRepository interface
 * and provides an access to the postgres database
 */
type ReviewRepository struct {
	db *postgres.DB
}

// NewReviewRepository creates a new review repository instance
func NewReviewRepository(db *postgres.DB) *ReviewRepository {
	return &ReviewRepository{
		db,
	}
}

// CreateReview creates a new review in the database
func (rp *ReviewRepository) CreateReview(ctx context.Context, review *review.ReviewDocument) (*review.ReviewDocument, error) {
	query := rp.db.QueryBuilder.Insert("reviews").
		Columns(
			"gig_id",
			"reviewer_id",
			"order_id",
			"seller_id",
			"review",
			"reviewer_image",
			"reviewer_username",
			"country",
			"review_type",
			"rating",
			"created_at",
		).
		Values(
			review.GigId,
			review.ReviewerId,
			review.OrderId,
			review.SellerId,
			review.Review,
			review.ReviewerImage,
			review.ReviewerUsername,
			review.Country,
			review.ReviewType,
			review.Rating,
			review.CreatedAt,
		).Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = rp.db.QueryRow(ctx, sql, args...).Scan(
		&review.Id,
		&review.GigId,
		&review.ReviewerId,
		&review.OrderId,
		&review.SellerId,
		&review.Review,
		&review.ReviewerImage,
		&review.ReviewerUsername,
		&review.Country,
		&review.ReviewType,
		&review.Rating,
		&review.CreatedAt,
	)
	if err != nil {
		if errCode := rp.db.ErrorCode(err); errCode == "23505" {
			return nil, postgres.ErrConflictingData
		}

		return nil, err
	}

	return review, nil
}

// GetReviewsByGigId gets all reviews in the database for a gig
func (rp *ReviewRepository) GetReviewsByGigId(ctx context.Context, id string) ([]*review.ReviewDocument, error) {
	reviews := make([]*review.ReviewDocument, 0)

	query := rp.db.QueryBuilder.Select("*").
		From("reviews").
		Where(sq.Eq{"gig_id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := rp.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var review review.ReviewDocument

		err := rows.Scan(
			&review.Id,
			&review.GigId,
			&review.ReviewerId,
			&review.OrderId,
			&review.SellerId,
			&review.Review,
			&review.ReviewerImage,
			&review.ReviewerUsername,
			&review.Country,
			&review.ReviewType,
			&review.Rating,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, &review)
	}

	return reviews, nil
}

// GetReviewsBySellerId gets all reviews in the database for a seller
func (rp *ReviewRepository) GetReviewsBySellerId(ctx context.Context, id string) ([]*review.ReviewDocument, error) {
	reviews := make([]*review.ReviewDocument, 0)

	query := rp.db.QueryBuilder.Select("*").
		From("reviews").
		Where(sq.Eq{"seller_id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := rp.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var review review.ReviewDocument

		err := rows.Scan(
			&review.Id,
			&review.GigId,
			&review.ReviewerId,
			&review.OrderId,
			&review.SellerId,
			&review.Review,
			&review.ReviewerImage,
			&review.ReviewerUsername,
			&review.Country,
			&review.ReviewType,
			&review.Rating,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, &review)
	}

	return reviews, nil
}
