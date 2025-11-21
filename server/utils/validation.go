package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	_required string = "required"
	_min      string = "min"
	_max      string = "max"
	_email    string = "email"
	_phone    string = "phone"

	ONE_MB = 1 << 20
)

func formatValidatorError(valErrors validator.ValidationErrors) string {
	allErrors := []string{}
	for _, validationErr := range valErrors {
		tagName := validationErr.Tag()
		fieldName := validationErr.Field()
		validationParam := validationErr.Param()

		errorMsg := ""
		switch tagName {
		case _required:
			errorMsg = fmt.Sprintf("%s %s", fieldName, "is required")
		case _min:
			errorMsg = fmt.Sprintf("%s %s %s %s", fieldName, "length must be greater than", validationParam, "characters")
		case _max:
			errorMsg = fmt.Sprintf("%s %s %s %s", fieldName, "length must be less than", validationParam, "characters")
		case _email:
			errorMsg = fmt.Sprintf("%s %s", fieldName, "must be a valid email address")
		case _phone:
			errorMsg = fmt.Sprintf("%s %s", fieldName, "must be a valid phone number")
		default:
			errorMsg = fmt.Sprintf("%s %s", fieldName, "is invalid")
		}

		allErrors = append(allErrors, errorMsg)
	}

	return strings.Join(allErrors, "\n")
}

func ValidateData[T any](validated *T) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(validated); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.New(formatValidatorError(validationErrors))
		}
		return errors.New("invalid data provided")
	}

	return nil
}

func ParseBodyAndValidate[T any](w http.ResponseWriter, r *http.Request, validated *T) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return errors.New("only application/json Content-Type is supported")
	}

	r.Body = http.MaxBytesReader(w, r.Body, ONE_MB) // Limit body size to 1MB
	if err := json.NewDecoder(r.Body).Decode(&validated); err != nil {
		http.Error(w, "Error Parsing JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}

	return ValidateData(validated)
}

func ParseInt32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}
