package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/generated"
	"github.com/linkc0829/go-ics/internal/graph/models"
)


func (r *queryResolver) MyPortfolio(ctx context.Context) (*models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyHistory(ctx context.Context, rangeArg int) (*models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetPortfolio(ctx context.Context, userID string) (*models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetHistory(ctx context.Context, userID string, rangeArg int) (*models.Portfolio, error) {
	panic("not implemented")
}

