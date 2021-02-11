package resolvers

import (
	"context"

	"github.com/linkc0829/go-icsharing/internal/graph/models"
)

func (r *mutationResolver) CreateIncome(ctx context.Context, input models.CreateIncomeInput) (*models.Income, error) {
	income, err := r.CreatePortfolio(ctx, input, "income")
	if err != nil {
		return nil, err
	}
	result := (*income).(models.Income)
	return &result, nil
}

func (r *mutationResolver) UpdateIncome(ctx context.Context, id string, input models.UpdateIncomeInput) (*models.Income, error) {
	income, err := r.UpdatePortfolio(ctx, id, input, "income")
	if err != nil {
		return nil, err
	}
	ret := (*income).(models.Income)
	return &ret, nil
}

func (r *mutationResolver) DeleteIncome(ctx context.Context, id string) (bool, error) {
	return r.DeletePortfolio(ctx, id, "income")
}

func (r *mutationResolver) VoteIncome(ctx context.Context, id string) (int, error) {
	return r.VotePortfolio(ctx, id, "income")
}

type incomeResolver struct{ *Resolver }

func (r *incomeResolver) Vote(ctx context.Context, obj *models.Income) ([]*models.User, error) {
	return r.resolveUsers(ctx, obj.Vote...)
}

func (r *incomeResolver) Owner(ctx context.Context, obj *models.Income) (*models.User, error) {
	owner, err := r.resolveUsers(ctx, obj.Owner)
	if err != nil {
		return nil, err
	}
	return owner[0], nil
}
