package command

import (
	"log/slog"
	"reflect"
	"strings"
)

// MaskSlogValue receives any struct and returns a slog.GroupValue
// with the fields mapped from the json tag and masked according to the sensitive tag.
func MaskSlogValue[T any](input T) slog.Value {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return slog.AnyValue(input)
	}

	t := v.Type()
	attrs := make([]slog.Attr, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Key name: uses the json tag, or the field name if there is none
		jsonTag := fieldType.Tag.Get("json")
		key := strings.Split(jsonTag, ",")[0]
		if key == "" || key == "-" {
			key = fieldType.Name
		}

		// Checks the sensitive tag
		sensitiveTag := fieldType.Tag.Get("sensitive")
		var value slog.Value
		if sensitiveTag == "true" {
			value = slog.StringValue("*************")
		} else {
			value = slog.AnyValue(field.Interface())
		}

		attrs = append(attrs, slog.Attr{Key: key, Value: value})
	}

	return slog.GroupValue(attrs...)
}
