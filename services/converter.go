package services

import (
	"github.com/boardware-cloud/core/services/model"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
	"github.com/chenyunda218/golambda"
)

func AccountListBackward(list common.List[core.Account]) common.List[model.Account] {
	return common.List[model.Account]{
		Data: golambda.Map(list.Data, func(_ int, account core.Account) model.Account {
			return AccountBackward(account)
		}),
		Pagination: list.Pagination,
	}
}

func AccountBackward(account core.Account) model.Account {
	return model.Account{
		Entity: account,
	}
}
