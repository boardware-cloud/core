package controllers

import (
	api "github.com/boardware-cloud/core-api"
	core "github.com/boardware-cloud/core/services"

	"github.com/boardware-cloud/common/utils"
)

func AccountBackward(account core.Account) api.Account {
	return api.Account{
		Id:    utils.UintToString(account.ID),
		Email: account.Email,
		Role:  api.Role(account.Role),
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

func ServiceBackward(service core.Service) api.Service {
	return api.Service{
		Id:          utils.UintToString(service.ID),
		Name:        service.Name,
		Title:       service.Title,
		Description: service.Description,
		Url:         service.Url,
		Type:        api.ServiceType(service.Type),
	}
}

func PaginationBackward(pagination core.Pagination) api.Pagination {
	return api.Pagination{
		Index: pagination.Index,
		Limit: pagination.Limit,
		Total: pagination.Total,
	}
}

func ServiceListBackward(serviceList core.List[core.Service]) api.ServiceList {
	pagination := PaginationBackward(serviceList.Pagination)
	var list api.ServiceList
	for _, v := range serviceList.Data {
		list.Data = append(list.Data, ServiceBackward(v))
	}
	list.Pagination = pagination
	return list
}
