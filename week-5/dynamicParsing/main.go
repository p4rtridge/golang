package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type Product struct {
	Name     string `json:"name"`
	Image    []byte `json:"image"`
	Quantity int    `json:"quantity"`
}

func populateStructFromForm(r *http.Request, result interface{}) error {
	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	// Get the value of the result (the struct) and its type
	v := reflect.ValueOf(result).Elem()
	t := v.Type()

	// Iterate through the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		formKey := field.Tag.Get("json") // Get the json tag as the form field name
		formValue := r.FormValue(formKey)

		// If no form value is present, skip this field
		if formValue == "" {
			continue
		}

		// Set the field based on its type
		switch field.Type.Kind() {
		case reflect.String:
			v.Field(i).SetString(formValue)
		case reflect.Int:
			intValue, err := strconv.Atoi(formValue)
			if err != nil {
				return fmt.Errorf("invalid int value for field %s: %v", formKey, err)
			}
			v.Field(i).SetInt(int64(intValue))
		case reflect.Float32:
			floatValue, err := strconv.ParseFloat(formValue, 32)
			if err != nil {
				return fmt.Errorf("invalid float value for field %s: %v", formKey, err)
			}
			v.Field(i).SetFloat(floatValue)
		default:
			return fmt.Errorf("unsupported field type %s", field.Type.Kind())
		}
	}
	return nil
}

func main() {
	product := Product{Name: "your mom", Image: make([]byte, 0), Quantity: 10}

	r := reflect.ValueOf(&product).Elem()

	for i := 0; i < r.NumField(); i++ {
		fieldType := r.Field(i).Type()
		fieldKind := r.Field(i).Kind()

		fmt.Println(fieldType, fieldKind)
	}
}
