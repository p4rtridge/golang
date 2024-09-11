package validate

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	return validator.New()
}

func ParseJsonField(model interface{}, structField string) string {
	t := reflect.TypeOf(model)
	field, found := t.FieldByName(structField)
	if found && field.Tag.Get("json") != "-" {
		return field.Tag.Get("json")
	} else {
		return structField
	}
}
