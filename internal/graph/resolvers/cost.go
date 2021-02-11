package resolvers

import (
	"context"

	"github.com/linkc0829/go-icsharing/internal/graph/models"
)

func (r *mutationResolver) CreateCost(ctx context.Context, input models.CreateCostInput) (*models.Cost, error) {
	cost, err := r.CreatePortfolio(ctx, input, "cost")
	if err != nil {
		return nil, err
	}
	result := (*cost).(models.Cost)
	return &result, nil
}

func (r *mutationResolver) UpdateCost(ctx context.Context, id string, input models.UpdateCostInput) (*models.Cost, error) {
	cost, err := r.UpdatePortfolio(ctx, id, input, "cost")
	if err != nil {
		return nil, err
	}
	ret := (*cost).(models.Cost)
	return &ret, nil
}

func (r *mutationResolver) DeleteCost(ctx context.Context, id string) (bool, error) {
	return r.DeletePortfolio(ctx, id, "cost")
}

func (r *mutationResolver) VoteCost(ctx context.Context, id string) (int, error) {
	return r.VotePortfolio(ctx, id, "cost")
}

type costResolver struct{ *Resolver }

func (r *costResolver) Vote(ctx context.Context, obj *models.Cost) ([]*models.User, error) {
	return r.resolveUsers(ctx, obj.Vote...)
}

func (r *costResolver) Owner(ctx context.Context, obj *models.Cost) (*models.User, error) {
	owner, err := r.resolveUsers(ctx, obj.Owner)
	if err != nil {
		return nil, err
	}
	return owner[0], nil
}
