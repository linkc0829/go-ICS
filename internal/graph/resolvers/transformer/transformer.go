package transformer

import (
	dbModel "github.com/linkc0829/go-icsharing/internal/db/mongodb/models"
	gqlModel "github.com/linkc0829/go-icsharing/internal/graph/models"
)

func DBPortfolioToGQLPortfolio(db dbModel.PortfolioModel, pType string) gqlModel.Portfolio {

	var gql gqlModel.Portfolio
	vote := []string{}
	for _, v := range db.Vote {
		vote = append(vote, v.Hex())
	}
	if pType == "cost" {
		gql = gqlModel.Cost{
			ID:          db.ID.Hex(),
			Owner:       db.Owner.Hex(),
			Amount:      db.Amount,
			Category:    db.Category,
			Description: db.Description,
			OccurDate:   db.OccurDate,
			Vote:        vote,
			Privacy:     db.Privacy,
		}
		return gql
	}
	gql = gqlModel.Income{
		ID:          db.ID.Hex(),
		Owner:       db.Owner.Hex(),
		Amount:      db.Amount,
		Category:    db.Category,
		Description: db.Description,
		OccurDate:   db.OccurDate,
		Vote:        vote,
		Privacy:     db.Privacy,
	}
	return gql
}
