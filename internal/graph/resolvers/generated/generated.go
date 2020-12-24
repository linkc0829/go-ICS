package generated

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/generated"
	"github.com/linkc0829/go-ics/internal/graph/models"
)

type Resolver struct{}

func (r *costResolver) Owner(ctx context.Context, obj *models.Cost) (*models.User, error) {
	panic("not implemented")
}

func (r *costResolver) Vote(ctx context.Context, obj *models.Cost) ([]*models.User, error) {
	panic("not implemented")
}

func (r *incomeResolver) Owner(ctx context.Context, obj *models.Income) (*models.User, error) {
	panic("not implemented")
}

func (r *incomeResolver) Vote(ctx context.Context, obj *models.Income) ([]*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input models.UpdateUserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateIncome(ctx context.Context, input models.CreateIncomeInput) (*models.Income, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateIncome(ctx context.Context, id string, input models.UpdateIncomeInput) (*models.Income, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteIncome(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateCost(ctx context.Context, input models.CreateCostInput) (*models.Cost, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateCost(ctx context.Context, id string, input models.UpdateCostInput) (*models.Cost, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteCost(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddFriend(ctx context.Context, id string) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) VoteCost(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *mutationResolver) VoteIncome(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) MyCostHistory(ctx context.Context, rangeArg int) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyIncomeHistory(ctx context.Context, rangeArg int) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyIncome(ctx context.Context) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyCost(ctx context.Context) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) MyFriends(ctx context.Context) ([]*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserIncome(ctx context.Context, id string) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserCost(ctx context.Context, id string) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserIncomeHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUserCostHistory(ctx context.Context, id string, rangeArg int) ([]models.Portfolio, error) {
	panic("not implemented")
}

func (r *userResolver) Friends(ctx context.Context, obj *models.User) ([]*models.User, error) {
	panic("not implemented")
}

func (r *userResolver) Followers(ctx context.Context, obj *models.User) ([]*models.User, error) {
	panic("not implemented")
}

// Cost returns generated.CostResolver implementation.
func (r *Resolver) Cost() generated.CostResolver { return &costResolver{r} }

// Income returns generated.IncomeResolver implementation.
func (r *Resolver) Income() generated.IncomeResolver { return &incomeResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type costResolver struct{ *Resolver }
type incomeResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
