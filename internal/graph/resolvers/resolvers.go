package resolvers

import (
	"context"

	"github.com/linkc0829/go-icsharing/internal/db/mongodb"
	"github.com/linkc0829/go-icsharing/internal/graph/generated"
	"github.com/linkc0829/go-icsharing/internal/graph/models"
)

//Resolver contains db element
type Resolver struct {
	DB *mongodb.MongoDB
}

// Cost returns generated.CostResolver implementation.
func (r *Resolver) Cost() generated.CostResolver { return &costResolver{r} }

// Income returns generated.IncomeResolver implementation.
func (r *Resolver) Income() generated.IncomeResolver { return &incomeResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *Resolver) resolveUsers(ctx context.Context, ids ...string) ([]*models.User, error) {
	result := make([]*models.User, len(ids))
	for i, id := range ids {
		user, err := r.Query().GetUser(ctx, id)
		if err != nil {
			return nil, err
		}
		result[i] = user
	}
	return result, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
