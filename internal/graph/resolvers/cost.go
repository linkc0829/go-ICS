package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/models"
)

func (r *mutationResolver) CreateCost(ctx context.Context, input models.CostInput) (*models.Cost, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateCost(ctx context.Context, id string, input models.CostInput) (*models.Cost, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteCost(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) LikeCost(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

type costResolver struct{ *Resolver }

func (r *costResolver) Vote(ctx context.Context, obj *models.Cost) ([]*models.User, error) {
	panic("not implemented")
}
