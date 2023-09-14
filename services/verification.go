package services

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"

	errorCode "github.com/boardware-cloud/common/code"
	constants "github.com/boardware-cloud/common/constants/account"
	model "github.com/boardware-cloud/model/core"
)

func GetVerification(identity string, purpose constants.VerificationCodePurpose) *model.VerificationCode {
	var verificationCode model.VerificationCode
	ctx := DB.Where("identity = ? AND purpose = ?",
		identity,
		purpose,
	).Order("created_at DESC").Find(&verificationCode)
	if ctx.RowsAffected == 0 {
		return nil
	}
	return &verificationCode
}

func CreateVerificationCode(identity string, purpose constants.VerificationCodePurpose) error {
	ctx := DB.Where("email = ?", identity).Find(&model.Account{})
	if purpose == constants.CREATE_ACCOUNT && ctx.RowsAffected != 0 {
		return errorCode.ErrEmailExists
	}
	if purpose == constants.SET_PASSWORD && ctx.RowsAffected == 0 {
		return errorCode.ErrNotFound
	}
	var verificationCode model.VerificationCode
	ctx = DB.Where("identity = ? AND purpose = ?",
		identity,
		purpose,
	).Order("created_at DESC").Find(&verificationCode)
	if ctx.RowsAffected == 0 || time.Now().Unix()-verificationCode.CreatedAt.Unix() >= 60 {
		newCode := &model.VerificationCode{
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
		if err != nil {
			return errorCode.ErrUndefined
		}
		DB.Delete(&model.VerificationCode{Identity: identity, Purpose: purpose})
		DB.Save(&newCode)
		return nil
	}
	return errorCode.ErrTooManyRequests
}

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
