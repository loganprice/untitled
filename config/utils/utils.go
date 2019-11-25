package configutils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	//errInvalidValue returned when the value passed to unsmarshl is nil or not a pointer to a struct
	errInvalidValue = errors.New("value must be a non-nil pointer to a struct or map")

	// errUnsupportedType returned when a field with tag "env" is unsupported.
	errUnsupportedType = errors.New("field is an unsupported type")

	// errUnexportedField returned when a field with tag "env" is not exported.
	errUnexportedField = errors.New("field must be exported")
)

// ParseKVStringSlice ...
func ParseKVStringSlice(input []string, seperator string) map[string]string {
	valueMap := make(map[string]string)
	for _, kv := range input {
		kvSplit := strings.SplitN(kv, seperator, 2)
		valueMap[kvSplit[0]] = kvSplit[1]
	}
	return valueMap
}

// Unmarshal ...
func Unmarshal(es map[string]string, tag string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errInvalidValue
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errInvalidValue
	}

	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		valueField := rv.Field(i)
		switch valueField.Kind() {
		case reflect.Struct:
			if !valueField.Addr().CanInterface() {
				continue
			}

			iface := valueField.Addr().Interface()
			err := Unmarshal(es, tag, iface)
			if err != nil {
				return err
			}
		}

		typeField := t.Field(i)
		value := typeField.Tag.Get(tag)
		if value == "" {
			continue
		}

		if !valueField.CanSet() {
			return errUnexportedField
		}

		envVar, ok := es[value]
		if !ok {
			continue
		}

		err := set(typeField.Type, valueField, envVar)
		if err != nil {
			return err
		}
		delete(es, tag)
	}

	return nil
}

func set(t reflect.Type, f reflect.Value, value string) error {
	switch t.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(t.Elem())
		err := set(t.Elem(), ptr.Elem(), value)
		if err != nil {
			return err
		}
		f.Set(ptr)
	case reflect.String:
		f.SetString(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		f.SetBool(v)
	case reflect.Int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		f.SetInt(int64(v))
	default:
		return errUnsupportedType
	}

	return nil
}
