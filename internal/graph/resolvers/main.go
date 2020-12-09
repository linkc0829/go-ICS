package resolvers

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/generated"
	"github.com/linkc0829/go-ics/internal/graph/models"
	"github.com/linkc0829/go-ics/internal/mongodb"
)

type Resolver struct{
	DB mongo.MongoDB
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, userID string) (bool, error) {
	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context, userID *string) ([]*models.User, error) {

	id := "ec17af15-e354-440c-a09f-69715fc8b595"
	email := "your@gmail.com"
	tmpuserID := "UserID-1"

    records := []*models.User{
        &models.User{
            ID:    	&id,
            Email:  &email,
            UserID: &tmpuserID,
		},
    }
    return records, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
