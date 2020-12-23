package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/models"
)

func (r *queryResolver) MyPortfolio(ctx context.Context) (*models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyHistory(ctx context.Context, rangeArg int) (*models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserPortfolio(ctx context.Context, userID string) (*models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserHistory(ctx context.Context, userID string, rangeArg int) (*models.Portfolio, error) {
	panic("not implemented")
}
