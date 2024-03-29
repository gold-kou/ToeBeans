package http

import (
	"fmt"
	"regexp"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	MinPasswordLength = 6
	MinVarcharLength  = 2
	MaxVarcharLength  = 255
	UUIDLength        = 36

	/* #nosec */
	errMsgPasswordValidation = "Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters"
)

var IsAlphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func (req *RequestRegisterUser) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&req.Email, validation.Required, is.Email, validation.Length(MinVarcharLength, MaxVarcharLength)),
		validation.Field(&req.Password, validation.Required, validation.By(PasswordValidation)))
	return validation.ValidateStruct(req, fieldRules...)
}

func (req *RequestLogin) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&req.Email, validation.Required, is.Email, validation.Length(MinVarcharLength, MaxVarcharLength)),
		validation.Field(&req.Password, validation.Required, validation.By(PasswordValidation)))
	return validation.ValidateStruct(req, fieldRules...)
}

func (req *RequestChangePassword) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&req.OldPassword, validation.Required, validation.By(PasswordValidation)),
		validation.Field(&req.NewPassword, validation.Required, validation.By(PasswordValidation)))
	return validation.ValidateStruct(req, fieldRules...)
}

func (req *RequestResetPassword) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&req.UserName, validation.Required, validation.Length(MinVarcharLength, MaxVarcharLength), is.Alphanumeric),
		validation.Field(&req.Password, validation.Required, validation.By(PasswordValidation)),
		validation.Field(&req.PasswordResetKey, validation.Required, is.UUID))
	return validation.ValidateStruct(req, fieldRules...)
}

func (req *RequestRegisterPosting) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&req.Title, validation.Required, validation.Length(MinVarcharLength, MaxVarcharLength)),
		validation.Field(&req.Image, validation.Required))
	return validation.ValidateStruct(req, fieldRules...)
}

func (e *RequestSendPasswordResetEmail) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&e.Email, validation.Required, is.Email, validation.Length(MinVarcharLength, MaxVarcharLength)))
	return validation.ValidateStruct(e, fieldRules...)
}

func (c *RequestRegisterComment) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&c.Comment, validation.Required, validation.Length(MinVarcharLength, MaxVarcharLength)))
	return validation.ValidateStruct(c, fieldRules...)
}

func (r *RequestSubmitUserReport) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&r.Detail, validation.Required))
	return validation.ValidateStruct(r, fieldRules...)
}

func (r *RequestSubmitPostingReport) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&r.Detail, validation.Required))
	return validation.ValidateStruct(r, fieldRules...)
}

// custom password validation
//
// upp: at least one upper case letter.
// low: at least one lower case letter.
// num: at least one digit.
// tot: at least eight characters long.
// No empty string or whitespace.
func PasswordValidation(pass interface{}) error {
	var (
		upp, low, num bool
		tot           uint8
	)

	for _, char := range pass.(string) {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		default:
			return fmt.Errorf(errMsgPasswordValidation)
		}
	}

	if !upp || !low || !num || tot < MinPasswordLength {
		return fmt.Errorf(errMsgPasswordValidation)
	}

	return nil
}
