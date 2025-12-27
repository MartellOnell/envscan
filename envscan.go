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
