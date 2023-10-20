package controllers

import (
	"github.com/boardware-cloud/common/utils"
	api "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"
	model "github.com/boardware-cloud/core/services"
)

func AccountBackward(account model.Account) api.Account {
	return api.Account{
		Id:           utils.UintToString(account.ID()),
		Email:        account.Email(),
		Role:         api.Role(account.Role()),
		HasTotp:      account.HasTotp(),
		RegisteredOn: account.RegisteredOn().Unix(),
	}
}

func SessionBackward(session core.Session) api.Session {
	return api.Session{
		ExpiredAt: session.ExpiredAt,
		Status:    api.SessionStatus(session.Status),
	}
}

func PaginationBackward(pagination core.Pagination) api.Pagination {
	return api.Pagination{
		Index: pagination.Index,
		Limit: pagination.Limit,
		Total: pagination.Total,
	}
}
