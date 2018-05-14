package forms

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/bluele/gforms"
)

type baseStructForm struct {
	innerform *gforms.FormInstance
}

func (r *baseStructForm) IsValid() bool {
	return r.innerform.IsValid()
}

func (r *baseStructForm) Errors() gforms.Errors {
	return r.innerform.Errors()
}

func (r *baseStructForm) CleanedData() map[string]interface{} {
	return r.innerform.CleanedData
}

func (r *baseStructForm) Data() gforms.Data {
	return r.innerform.Data
}

func (r *baseStructForm) Load(s interface{}) {
	r.innerform.MapTo(s)
}

// EmailForm ...
type EmailForm struct {
	baseStructForm

	Email string `form:"email"`
}

var emailGForm = gforms.DefineForm(
	gforms.NewFields(
		gforms.NewTextField(
			"Email",
			gforms.Validators{
				gforms.Required(),
				gforms.EmailValidator(),
				gforms.MaxLengthValidator(32),
			},
		),
	),
)

func MakeEmailForm(r *http.Request) *EmailForm {
	form := &EmailForm{}
	form.innerform = emailGForm(r)
	return form
}

type ResetPasswordForm struct {
	baseStructForm

	Password       string `form:"Password"`
	RepeatPassword string `form:"RepeatPassword"`
}

var resetPasswordGForm = gforms.DefineForm(
	gforms.NewFields(
		gforms.NewTextField(
			"Password",
			gforms.Validators{
				gforms.Required(),
				gforms.MinLengthValidator(6),
				gforms.MaxLengthValidator(20),
			},
		),
		gforms.NewTextField(
			"RepeatPassword",
			gforms.Validators{
				gforms.Required(),
				EqualToValidator("Password"),
			},
		),
	),
)

func MakeResetPasswordForm(r *http.Request) *ResetPasswordForm {
	form := &ResetPasswordForm{}
	form.innerform = resetPasswordGForm(r)
	return form
}

// LoginForm ...
type LoginForm struct {
	baseStructForm

	EmailOrPhone string `form:"email_or_phone"`
	Password     string `form:"password"`
	RememberMe   bool
}

var loginGForm = gforms.DefineForm(
	gforms.NewFields(
		gforms.NewTextField(
			"EmailOrPhone",
			gforms.Validators{
				gforms.Required(),
				emailOrPhoneValidator{
					Message: "Input Email Or Phone",
				}, // gforms.EmailValidator(),
				gforms.MaxLengthValidator(32),
			},
		),
		gforms.NewTextField(
			"Password",
			gforms.Validators{
				gforms.Required(),
				gforms.MinLengthValidator(6),
				gforms.MaxLengthValidator(20),
			},
		),
		gforms.NewBooleanField(
			"RememberMe",
			gforms.Validators{},
		),
	),
)

func MakeLoginForm(r *http.Request) *LoginForm {
	form := &LoginForm{}
	form.innerform = loginGForm(r)
	return form
}

type emailOrPhoneValidator struct {
	Message string
	gforms.Validator
}

func (vl emailOrPhoneValidator) Validate(fi *gforms.FieldInstance, fo *gforms.FormInstance) error {
	evl := gforms.EmailValidator()
	err := evl.Validate(fi, fo)
	if err != nil {
		goto PHONE
	} else {
		return nil
	}

PHONE:
	pvl := gforms.RegexpValidator(`^[\+\ \d]+$`)
	err = pvl.Validate(fi, fo)
	if err != nil {
		return errors.New(vl.Message)
	}

	return nil
}

type equalToValidator struct {
	Message string
	EqualTo string
	gforms.Validator
}

func (vl equalToValidator) Validate(fi *gforms.FieldInstance, fo *gforms.FormInstance) error {
	v := fi.V
	other, ok := fo.GetField(vl.EqualTo)
	if !ok {
		return errors.New("Can't Not Get Other Field")
	}
	vOther := other.GetV()
	if v.IsNil || v.Kind != reflect.String || v.Value == "" {
		if vOther.IsNil || vOther.Kind != reflect.String || vOther.Value == "" {
			return nil
		}
		return errors.New(vl.Message)
	}

	sv := v.Value.(string)
	svOther := vOther.Value.(string)

	if sv != svOther {
		return errors.New(vl.Message)
	}

	return nil
}

func EqualToValidator(equalto string, message ...string) equalToValidator {

	vl := equalToValidator{}
	vl.EqualTo = equalto
	if len(message) > 0 {
		vl.Message = message[0]
	} else {
		vl.Message = fmt.Sprintf("Ensure this value equal to %s", vl.EqualTo)
	}

	return vl
}
