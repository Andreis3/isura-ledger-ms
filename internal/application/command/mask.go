package command

import (
	"log/slog"
	"reflect"
	"strings"
)

// MaskSlogValue recebe qualquer struct e retorna um slog.GroupValue
// com os campos mapeados a partir da tag json e mascarados conforme a tag sensitive.
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

		// Nome da chave: usa a tag json, ou o nome do campo se não houver
		jsonTag := fieldType.Tag.Get("json")
		key := strings.Split(jsonTag, ",")[0]
		if key == "" || key == "-" {
			key = fieldType.Name
		}

		// Verifica a tag sensitive
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
