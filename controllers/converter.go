package controllers

import (
	api "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"

	"github.com/boardware-cloud/common/utils"
)

func AccountBackward(account core.Account) api.Account {
	return api.Account{
		Id:      utils.UintToString(account.ID),
		Email:   account.Email,
		Role:    api.Role(account.Role),
		HasTotp: account.HasTotp,
	}
}

func SessionBackward(session core.Session) api.Session {
	return api.Session{
		Account:     AccountBackward(session.Account),
		Token:       session.Token,
		TokenType:   string(session.TokeType),
		TokenFormat: string(session.TokenFormat),
		ExpiredAt:   session.ExpiredAt,
		Status:      api.SessionStatus(session.Status),
	}
}

func PaginationBackward(pagination core.Pagination) api.Pagination {
	return api.Pagination{
		Index: pagination.Index,
		Limit: pagination.Limit,
		Total: pagination.Total,
	}
}
