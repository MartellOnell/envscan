// Package envscan provides utilities for reading environment variables
// and populating struct fields using reflection and struct tags.
package envscan

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrNilPointerDeference = errors.New("nil pointer deference error")
	ErrVMustBePtr          = errors.New("v must be pointer on struct")
)

// ReadEnvironment reads environment variables and populates the fields of a struct.
//
// This function uses reflection to inspect struct fields and their 'env' tags,
// then assigns values from environment variables or default values to those fields.
//
// Parameters:
//   - v: A pointer to a struct whose fields will be populated. Must not be nil.
//   - defaultEnvData: A map of default values to use when environment variables are not set.
//     Keys should match the 'env' tag values.
//
// Requirements:
//   - v must be a non-nil pointer to a struct
//   - All struct fields must have an 'env' tag specifying the environment variable name
//   - Environment variables or default values must be set for all fields
//
// Supported field types:
//   - string: Direct assignment from environment variable
//   - bool: Parsed using strconv.ParseBool
//   - int, int8, int16, int32, int64: Parsed as base-10 integers
//   - []string: Comma-separated values split into a slice
//
// The function follows this priority for value assignment:
//  1. Environment variable value (if set and non-empty)
//  2. Default value from defaultEnvData map (if environment variable is empty)
//
// Example:
//
//	type Config struct {
//	    Host     string   `env:"APP_HOST"`
//	    Port     int      `env:"APP_PORT"`
//	    Debug    bool     `env:"APP_DEBUG"`
//	    Features []string `env:"APP_FEATURES"`
//	}
//
//	config := &Config{}
//	defaults := map[string]string{
//	    "APP_PORT": "8080",
//	    "APP_DEBUG": "false",
//	}
//
//	err := ReadEnvironment(config, defaults)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Returns:
//   - ErrNilPointerDeference if v is nil
//   - ErrVMustBePtr if v is not a pointer to a struct
//   - An error if any field is missing an 'env' tag
//   - An error if any environment variable is not set and has no default
//   - An error if type conversion fails for bool or integer fields
//   - An error if a field type is unsupported
func ReadEnvironment(v any, defaultEnvData map[string]string) error {
	if v == nil {
		return ErrNilPointerDeference
	}

	refVal := reflect.ValueOf(v)

	if refVal.Kind() == reflect.Ptr {
		refVal = reflect.Indirect(refVal)
	} else {
		return ErrVMustBePtr
	}

	if refVal.Kind() == reflect.Interface {
		refVal = refVal.Elem()
	}

	if refVal.Kind() != reflect.Struct {
		return ErrVMustBePtr
	}

	refType := reflect.TypeOf(refVal.Interface())

	for i := range refVal.NumField() {
		fieldVal := refVal.Field(i)
		fieldType := refType.Field(i)

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			return fmt.Errorf("Struct field \"%s\" is missing 'env' tag", fieldType.Name)
		}

		valueToAssign := os.Getenv(envTag)
		if valueToAssign == "" {
			valueToAssign = defaultEnvData[envTag]
		}

		if valueToAssign == "" {
			return fmt.Errorf("environment variable %s not set", envTag)
		}

		if !fieldVal.IsValid() || !fieldVal.CanAddr() || !fieldVal.CanSet() {
			return fmt.Errorf("cannot assign to field %s", fieldType.Name)
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(valueToAssign)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(valueToAssign)
			if err != nil {
				return fmt.Errorf("failed to parse bool for field %s: %w", fieldType.Name, err)
			}
			fieldVal.SetBool(boolVal)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(valueToAssign, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int for field %s: %w", fieldType.Name, err)
			}
			fieldVal.SetInt(intVal)
		case reflect.Array, reflect.Slice:
			// Currently only supports []string
			if fieldVal.Type().Elem().Kind() != reflect.String {
				return fmt.Errorf("unsupported slice element type %s for field %s", fieldVal.Type().Elem().Kind().String(), fieldType.Name)
			}
			strArr := strings.Split(valueToAssign, ",")
			fieldVal.Set(reflect.ValueOf(strArr))
		default:
			return fmt.Errorf("unsupported field type %s for field %s", fieldVal.Kind().String(), fieldType.Name)
		}
	}

	return nil
}
