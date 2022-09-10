package repository

import (
	"context"

	"gitlab.com/dh-backend/auth-service/internal/model"
)

func (r *Repository) CreateInterest(ctx context.Context, arg *model.Interest) error {
	return r.db.Save(arg)
}

func (r *Repository) CreateUserInterest(ctx context.Context, userId string, interestId string) error {
	arg := &model.UserInterest{
		UserId:     userId,
		InterestId: interestId,
	}

	return r.db.Save(arg)
}
