package enum

import "reflect"

func New[T any]() T {
	var enum T

	enumType := reflect.ValueOf(enum).Type()
	for i := range enumType.NumField() {
		field := enumType.Field(i)

		value := field.Name
		if tag, ok := field.Tag.Lookup("value"); ok {
			value = tag
		}

		fieldTo := reflect.ValueOf(&enum).Elem().FieldByName(field.Name)
		converted := reflect.ValueOf(value).Convert(fieldTo.Type())

		fieldTo.Set(converted)
	}

	return enum
}

func Validator[T any, V comparable]() func(value V) bool {
	var enum T
	enumType := reflect.ValueOf(enum).Type()

	validValues := make(map[V]bool)
	for i := range enumType.NumField() {
		field := enumType.Field(i)

		value := field.Name
		if tag, ok := field.Tag.Lookup("value"); ok {
			value = tag
		}

		fieldTo := reflect.ValueOf(&enum).Elem().FieldByName(field.Name)
		converted := reflect.ValueOf(value).Convert(fieldTo.Type())

		var enumValue V
		reflect.ValueOf(&enumValue).Elem().Set(converted)

		validValues[enumValue] = true
	}

	return func(value V) bool {
		if value, ok := validValues[value]; ok {
			return value
		}
		return false
	}
}
