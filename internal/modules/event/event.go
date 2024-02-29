package event

import (
	"context"
	wrapper "payment-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindEventById(ctx context.Context, id string) <-chan wrapper.Result
}
