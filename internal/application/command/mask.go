package command

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

// MaskSlogValue receives any struct and returns a slog.GroupValue
// with the fields mapped from the json tag and masked according to the sensitive tag.
//
// Options for the sensitive tag:
//   - sensitive:"true" or sensitive:"full": Masks the entire value ("*************")
//   - sensitive:"partial": Preserves the middle digits and masks the rest (e.g., "***45678***")
func MaskSlogValue[T any](input T) slog.Value {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return slog.AnyValue(nil)
		}
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

		// Key name: uses the json tag or the field name if there is no tag
		jsonTag := fieldType.Tag.Get("json")
		key := strings.Split(jsonTag, ",")[0]
		if key == "" || key == "-" {
			key = fieldType.Name
		}

		// Check the sensitive tag
		sensitiveTag := fieldType.Tag.Get("sensitive")
		var value slog.Value

		switch sensitiveTag {
		case "true", "full":
			value = slog.StringValue("*************")
		case "partial":
			strVal := fmt.Sprintf("%v", field.Interface())
			value = slog.StringValue(maskMiddleVisible(strVal))
		default:
			value = slog.AnyValue(field.Interface())
		}

		attrs = append(attrs, slog.Attr{Key: key, Value: value})
	}

	return slog.GroupValue(attrs...)
}

// maskMiddleVisible hides the beginning and end of the string, keeping the middle characters visible.
func maskMiddleVisible(s string) string {
	runes := []rune(s)
	length := len(runes)

	// If it's a very short string, mask the whole thing for security
	if length <= 4 {
		return "******"
	}

	// Defines the number of characters to be displayed in the middle (e.g., 1/3 of the size)
	visibleLen := length / 3
	if visibleLen < 2 {
		visibleLen = 2
	}

	start := (length - visibleLen) / 2
	end := start + visibleLen

	var sb strings.Builder
	for i := 0; i < start; i++ {
		sb.WriteRune('*')
	}
	sb.WriteString(string(runes[start:end]))
	for i := end; i < length; i++ {
		sb.WriteRune('*')
	}

	return sb.String()
}
