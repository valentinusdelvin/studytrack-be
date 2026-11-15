package validation

import (
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/valentinusdelvin/studytrack-be/domain/errx"
)

func Validate(value any) error {
	valueOf := reflect.ValueOf(value)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	typeOf := valueOf.Type()
	var messages []string
	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		tag := typeOf.Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}
		for _, rule := range strings.Split(tag, ",") {
			parts := strings.SplitN(rule, "=", 2)
			name, arg := parts[0], ""
			if len(parts) == 2 {
				arg = parts[1]
			}
			if !valid(field, name, arg) {
				messages = append(messages, fmt.Sprintf("%s failed %s validation", strings.ToLower(typeOf.Field(i).Name), rule))
			}
		}
	}
	if len(messages) > 0 {
		return errx.ErrBadRequest.WithMessage("The request body is not valid.").WithDetails(map[string]any{"validation": messages})
	}
	return nil
}

func valid(value reflect.Value, rule, arg string) bool {
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return rule != "required"
		}
		value = value.Elem()
	}
	text := fmt.Sprint(value.Interface())
	empty := strings.TrimSpace(text) == "" || (value.Kind() == reflect.Uint64 && value.Uint() == 0)
	switch rule {
	case "omitempty":
		return true
	case "required":
		return !empty
	case "email":
		address, err := mail.ParseAddress(text)
		return err == nil && address.Address == text
	case "min":
		limit, err := strconv.Atoi(arg)
		if err != nil {
			return false
		}
		if value.Kind() >= reflect.Int && value.Kind() <= reflect.Int64 {
			return value.Int() >= int64(limit)
		}
		if value.Kind() == reflect.Uint64 {
			return value.Uint() >= uint64(limit)
		}
		return utf8.RuneCountInString(text) >= limit
	case "max":
		limit, err := strconv.Atoi(arg)
		if err != nil {
			return false
		}
		return utf8.RuneCountInString(text) <= limit
	case "oneof":
		for _, option := range strings.Fields(arg) {
			if text == option {
				return true
			}
		}
		return false
	default:
		return true
	}
}
