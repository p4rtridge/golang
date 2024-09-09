package helpers

import "reflect"

func GetJsonField(model interface{}, structField string) string {
	t := reflect.TypeOf(model)
	field, found := t.FieldByName(structField)
	if found {
		return field.Tag.Get("json")
	} else {
		return structField
	}
}
