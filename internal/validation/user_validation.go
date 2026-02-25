package validation

import (
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	errs "github.com/StewardMcCormick/Paste_Bin/internal/error"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var (
	RequiredFieldMessage = "this field is required"
	MinLengthMessage     = "Minimum length - "
	MaxLengthMessage     = "Maximum length - "
)

type userValidator struct {
	valid *validator.Validate
}

func NewUserValidator(valid *validator.Validate) *userValidator {
	return &userValidator{valid: valid}
}

func (uv *userValidator) Validate(user *dto.UserRequest) error {
	if err := uv.valid.Struct(user); err != nil {
		return uv.mapValidErrorToCustomError(err.(validator.ValidationErrors))
	}

	return nil
}

func (uv *userValidator) mapValidErrorToCustomError(err validator.ValidationErrors) error {
	ve := errs.ValidationError{
		Message: errs.UserValidationError.Error(),
		Status:  http.StatusBadRequest,
		Errors:  make([]errs.ValidationFieldError, len(err)),
	}

	for i, fe := range err {
		ve.Errors[i] = errs.ValidationFieldError{
			Field:   fe.Field(),
			Message: uv.convertTagToMessage(fe),
			Value:   fe.Value(),
		}
	}

	return ve
}

func (uv *userValidator) convertTagToMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return RequiredFieldMessage
	case "min":
		return MinLengthMessage + fe.Param()
	case "max":
		return MaxLengthMessage + fe.Param()
	default:
		return fe.Error()
	}
}
