package services

import (
	"bytes"
	"fmt"
	"math/rand"
	"text/template"
	"time"

	"github.com/Dparty/common/singleton"
	errorCode "github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	coreModel "github.com/boardware-cloud/model/core"
)

const charset = "0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomNumberString(length int) string {
	return StringWithCharset(length, charset)
}

var verificationCodeService = singleton.NewSingleton[VerificationCodeService](newVerificationCodeService, singleton.Eager)

func GetVerificationCodeService() *VerificationCodeService {
	return verificationCodeService.Get()
}

func newVerificationCodeService() *VerificationCodeService {
	return &VerificationCodeService{verificationCodeRepository: coreModel.NewVerificationCodeRepository()}
}

type VerificationCodeService struct {
	verificationCodeRepository *coreModel.VerificationCodeRepository
}

func (v VerificationCodeService) CreateVerificationCode(identity string, purpose constants.VerificationCodePurpose) error {
	account := accountRepository.GetByEmail(identity)
	if purpose == constants.CREATE_ACCOUNT && account == nil {
		return errorCode.ErrEmailExists
	}

	if purpose == constants.SET_PASSWORD && account == nil {
		return errorCode.ErrNotFound
	}
	var verificationCode coreModel.VerificationCode
	ctx := DB.Where("identity = ? AND purpose = ?",
		identity,
		purpose,
	).Order("created_at DESC").Find(&verificationCode)
	if ctx.RowsAffected == 0 || time.Now().Unix()-verificationCode.CreatedAt.Unix() >= 60 {
		newCode := &coreModel.VerificationCode{
			Identity: identity,
			Purpose:  purpose,
			Code:     RandomNumberString(6),
		}
		tmpl := template.New("verification")
		tmpl.Parse(VerificationEmailTemplate)
		var verificationCodeMap map[string]string = make(map[string]string)
		verificationCodeMap["VerificationCode"] = newCode.Code
		var htmlString bytes.Buffer
		tmpl.Execute(&htmlString, &verificationCodeMap)
		err := emailSender.SendHtml("", "Boardware Cloud verification code",
			htmlString.String(), []string{identity}, []string{}, []string{})
		fmt.Println(err)
		if err != nil {
			return errorCode.ErrUndefined
		}
		code := &coreModel.VerificationCode{Identity: identity, Purpose: purpose}
		DB.Where(code).Delete(code)
		DB.Save(&newCode)
		return nil
	}
	return errorCode.ErrTooManyRequests
}
