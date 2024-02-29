package user

import (
	"context"
	wrapper "payment-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindOneUserId(ctx context.Context, userId string) <-chan wrapper.Result
}
