package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func ParseQueryAndValidate(r *http.Request, dest any) error {
	if err := parseQueries(r, dest); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(dest); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.New(formatValidatorError(validationErrors))
		}
		return errors.New("invalid data provided")
	}
	return nil
}

func parseQueries(r *http.Request, dest any) error {
	queryParams := r.URL.Query()

	// Check type of destination
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	destElem := destValue.Elem()
	switch destElem.Kind() {
	case reflect.Struct:
		// For structs, we need to map query params to struct fields
		return parseIntoStruct(queryParams, destElem)
	case reflect.Map:
		// For maps, check if it's map[string]interface{}
		if destElem.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("map must have string keys")
		}
		return parseIntoMap(queryParams, destElem)
	default:
		return fmt.Errorf("destination must be a pointer to struct or map[string]interface{}")
	}
}

func parseIntoStruct(queryParams map[string][]string, destStruct reflect.Value) error {
	destType := destStruct.Type()

	// Iterate through all fields of the struct
	for i := range destStruct.NumField() {
		field := destStruct.Field(i)
		fieldType := destType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		if fieldType.Anonymous && field.Kind() == reflect.Struct {
			// Recursively parse the embedded struct
			if err := parseIntoStruct(queryParams, field); err != nil {
				return err
			}
			continue
		}

		// Get the query tag or use field name as fallback
		queryTag := fieldType.Tag.Get("query")
		if queryTag == "" {
			queryTag = strings.ToLower(fieldType.Name)
		} else if queryTag == "-" {
			// Skip this field if tag is "-"
			continue
		}

		// Check if param exists
		values, exists := queryParams[queryTag]
		if !exists || len(values) == 0 {
			continue // Skip if parameter not provided
		}

		// Get the first value (most URL params just have one value)
		value := values[0]

		// Set the field based on its type
		if err := setFieldValue(field, value); err != nil {
			return fmt.Errorf("error setting field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func parseIntoMap(queryParams map[string][]string, destMap reflect.Value) error {
	// Create a new map if the map is nil
	if destMap.IsNil() {
		destMap.Set(reflect.MakeMap(destMap.Type()))
	}

	// Add each query parameter to the map
	for key, values := range queryParams {
		// Handle multiple values
		if len(values) > 1 {
			destMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(values))
		} else if len(values) == 1 {
			// For single values, try to convert to appropriate type
			destMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(guessValueType(values[0])))
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(val)

	case reflect.Int32:
		val, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(val)

	case reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(val)

	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(val)

	case reflect.Struct:
		// Special handling for time.Time
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				// Try alternative formats
				t, err = time.Parse("2006-01-02", value)
				if err != nil {
					return err
				}
			}
			field.Set(reflect.ValueOf(t))
		} else {
			return fmt.Errorf("unsupported struct type")
		}

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

func guessValueType(value string) interface{} {
	// Try boolean
	if value == "true" || value == "false" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}

	// Try integer
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}

	// Try float
	if f, err := strconv.ParseFloat(value, 64); err == nil && !strings.Contains(value, ".") {
		return f
	}

	// Default to string
	return value
}

func ToJson(data any) []byte {
	val, _ := json.Marshal(data)
	return val
}
