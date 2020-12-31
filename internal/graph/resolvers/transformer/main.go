package transformer

import (
	gqlModel "github.com/linkc0829/go-ics/internal/graph/models"
	dbModel "github.com/linkc0829/go-ics/internal/mongodb/models"
)

func DBPortfolioToGQLPortfolio(db dbModel.PortfolioModel) gqlModel.Portfolio {

	var gql gqlModel.Portfolio
	switch m := db.(type) {
	case dbModel.CostModel:
		vote := []string{}
		for _, v := range m.Vote {
			vote = append(vote, v.Hex())
		}
		gql = gqlModel.Cost{
			ID:          m.ID.Hex(),
			Owner:       m.Owner.Hex(),
			Amount:      m.Amount,
			Category:    m.Category,
			Description: m.Description,
			OccurDate:   m.OccurDate,
			Vote:        vote,
		}
		return gql

	case dbModel.IncomeModel:
		vote := []string{}
		for _, v := range m.Vote {
			vote = append(vote, v.Hex())
		}
		gql = gqlModel.Income{
			ID:          m.ID.Hex(),
			Owner:       m.Owner.Hex(),
			Amount:      m.Amount,
			Category:    m.Category,
			Description: m.Description,
			OccurDate:   m.OccurDate,
			Vote:        vote,
		}
		return gql
	}
	return nil
}
