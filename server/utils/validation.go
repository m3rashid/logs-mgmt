package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
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

// ParseTimestamp parses a timestamp string into pgtype.Timestamptz
// Supports RFC3339, ISO8601, and common date formats
func ParseTimestamp(s string) (pgtype.Timestamptz, error) {
	var result pgtype.Timestamptz

	if s == "" {
		return result, nil
	}

	// Try multiple timestamp formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, s)
		if err == nil {
			result = pgtype.Timestamptz{Time: t, Valid: true}
			return result, nil
		}
	}

	return result, fmt.Errorf("unable to parse timestamp: %s", s)
}
