package services

import (
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
	"github.com/chenyunda218/golambda"
)

func AccountListBackward(list common.List[core.Account]) common.List[Account] {
	return common.List[Account]{
		Data: golambda.Map(list.Data, func(_ int, account core.Account) Account {
			return Account{Entity: account}
		}),
		Pagination: list.Pagination,
	}
}
