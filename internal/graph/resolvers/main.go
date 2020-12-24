package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/generated"
	"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/internal/mongodb"
)

//Resolver contains db element
type Resolver struct {
	DB *mongodb.MongoDB
}

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

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
