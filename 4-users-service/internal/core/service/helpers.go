package service

import (
	"fmt"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateErrorCheck(resp *mongo.UpdateResult, err error) error {
	if err != nil {
		return svc.NewError(svc.ErrInternalFailure, err)
	}
	if resp.MatchedCount == 0 {
		return svc.NewError(svc.ErrNotFound, fmt.Errorf("seller id does not exist"))
	}
	if resp.ModifiedCount == 0 {
		return svc.NewError(svc.ErrBadRequest, fmt.Errorf("seller was not updated"))
	}

	return nil
}
