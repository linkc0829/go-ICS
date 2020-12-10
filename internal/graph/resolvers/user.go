package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/generated"
	"github.com/linkc0829/go-ics/internal/graph/models"
)

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUser(ctx context.Context, userID string) (*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) MyFriends(ctx context.Context) (*models.Users, error) {
	panic("not implemented")
}

func (r *queryResolver) MyFollowers(ctx context.Context) (*models.Users, error) {
	panic("not implemented")
}


func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) AddFriends(ctx context.Context, userID string) (*models.User, error) {
	panic("not implemented")
}
