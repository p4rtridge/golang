package pkg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunInTransaction(ctx context.Context, db *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)

	if err == nil {
		return tx.Commit(ctx)
	}

	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}

func MultipartParser(m *multipart.Form, result interface{}) error {
	r := reflect.ValueOf(result)
	if r.Kind() != reflect.Pointer {
		return errors.New("expected pointer")
	}

	r = r.Elem()
	t := r.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		formKey, ok := field.Tag.Lookup("json")
		if !ok {
			return errors.New("must have json tag field")
		}

		// check if the field is []byte
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Uint8 {
			file, valid := m.File[formKey]
			if !valid && formKey != "-" {
				return errors.New("cannot found the target field")
			}
			// ignore json:"-"
			if len(file) == 0 {
				continue
			}
			f := file[0]

			fileContent, err := f.Open()
			if err != nil {
				return fmt.Errorf("error opening file for field %s: %v", formKey, err)
			}
			defer fileContent.Close()

			buff := bytes.NewBuffer(nil)
			if _, err := io.Copy(buff, fileContent); err != nil {
				return err
			}

			r.Field(i).SetBytes(buff.Bytes())
			continue
		}

		formValue := m.Value[formKey]
		if len(formValue) == 0 {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			r.Field(i).SetString(formValue[0])
		case reflect.Int:
			intValue, err := strconv.Atoi(formValue[0])
			if err != nil {
				return err
			}
			r.Field(i).SetInt(int64(intValue))
		case reflect.Float32:
			floatValue, err := strconv.ParseFloat(formValue[0], 32)
			if err != nil {
				return err
			}
			r.Field(i).SetFloat(floatValue)
		default:
			return fmt.Errorf("unsupported field type %s", field.Type.Kind())
		}
	}

	return nil
}
