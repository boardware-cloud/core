package services

import (
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"
	coreModel "github.com/boardware-cloud/model/core"
)

var DB = model.GetDB()

var emailSender = notifications.GetEmailSender()
var accountRepository = coreModel.GetAccountRepository()
var verificationCodeRepository = coreModel.GetVerificationCodeRepository()
var ticketRepository = coreModel.GetTicketRepository()
var webauthRepository = coreModel.GetWebauthRepository()
var sessionDataRepository = coreModel.GetSessionDataRepository()
