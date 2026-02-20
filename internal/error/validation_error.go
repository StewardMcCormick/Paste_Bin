package error

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StewardMcCormick/Paste_Bin/internal/util"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ValidationError struct {
	Message string                 `json:"message"`
	Status  int                    `json:"status"`
	Errors  []ValidationFieldError `json:"errors"`
}

type ValidationFieldError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

func SendValidationError(ctx context.Context, w http.ResponseWriter, status int, err validator.ValidationErrors) {
	log := util.GetLoggerFromCtx(ctx)

	w.WriteHeader(status)
	ve := ValidationError{
		Message: "Validation error",
		Status:  status,
		Errors:  make([]ValidationFieldError, len(err)),
	}

	for i, fe := range err {
		ve.Errors[i] = ValidationFieldError{
			Field:   fe.Field(),
			Message: convertTagToMessage(fe),
			Value:   fe.Value(),
		}
	}

	e := json.NewEncoder(w).Encode(ve)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(e.Error())
	}
	log.Info(fmt.Sprintf("validation error: %v", ve))
}

func convertTagToMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return fmt.Sprintf("Minimum lenght - %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum lenght - %s", fe.Param())
	default:
		return fe.Error()
	}
}
