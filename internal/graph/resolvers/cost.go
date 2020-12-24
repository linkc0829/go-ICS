package resolvers

import (
	"context"

	"github.com/linkc0829/go-ics/internal/graph/models"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
	"github.com/linkc0829/go-ics/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateCost(ctx context.Context, input models.CostInput) (*models.Cost, error) {
	me := ctx.Value(utils.ProjectContextKeys.UserCtxKey).(*dbModel.UserModel)

	newCost := &dbModel.CostModel{
		ID:     primitive.NewObjectID(),
		Owner:  me.ID,
		Amount: *input.Amount,
	}

	panic("not implemented")
}

func (r *mutationResolver) UpdateCost(ctx context.Context, id string, input models.CostInput) (*models.Cost, error) {
	panic("not implemented")
}

func (r *mutationResolver) DeleteCost(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) VoteCost(ctx context.Context, id string) (int, error) {
	panic("not implemented")
}

type costResolver struct{ *Resolver }

func (r *costResolver) Vote(ctx context.Context, obj *models.Cost) ([]*models.User, error) {
	panic("not implemented")
}

func (r *costResolver) Owner(ctx context.Context, obj *models.Cost) (*models.User, error) {
	panic("not implemented")
}
