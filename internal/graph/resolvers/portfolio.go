package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/models"
)

func (r *queryResolver) MyCostHistory(ctx context.Context) (models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyIncomeHistory(ctx context.Context) (models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyIncome(ctx context.Context, rangeArg int) (models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyCost(ctx context.Context, rangeArg int) (models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserCost(ctx context.Context, id string) (models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserIncome(ctx context.Context, id string) (models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserIncomeHistory(ctx context.Context, id string, rangeArg int) (models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserCostHistory(ctx context.Context, id string, rangeArg int) (models.Portfolio, error) {
	panic("not implemented")
}
