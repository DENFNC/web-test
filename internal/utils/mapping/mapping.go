package mapping

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapStructModel(src, dst any) error {
	srcVal, err := validateSrc(src)
	if err != nil {
		return err
	}

	dstVal, err := validateDst(dst)
	if err != nil {
		return err
	}

	srcType := srcVal.Type()

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		dstField, ok := dstVal.Type().FieldByName(srcField.Name)
		if !ok {
			continue
		}

		srcValue := srcVal.Field(i)
		dstValue := dstVal.FieldByName(srcField.Name)

		if !dstValue.CanSet() {
			continue
		}

		switch dstField.Type {
		case reflect.TypeOf(pgtype.Text{}):
			val := srcValue.Interface().(string)
			dstValue.Set(reflect.ValueOf(toPgText(val)))

		case reflect.TypeOf(pgtype.UUID{}):
			val := srcValue.Interface().(string)
			dstValue.Set(reflect.ValueOf(toPgUUID(val)))

		case reflect.TypeOf(pgtype.Timestamp{}):
			val := srcValue.Interface().(time.Time)
			dstValue.Set(reflect.ValueOf(toPgTimestamp(val)))

		case reflect.TypeOf(pgtype.Bool{}):
			val := srcValue.Interface().(bool)
			dstValue.Set(reflect.ValueOf(toPgBool(val)))

		case reflect.TypeOf(pgtype.Int2{}):
			val := srcValue.Interface().(int16)
			dstValue.Set(reflect.ValueOf(toPgInt16(val)))

		case reflect.TypeOf(pgtype.Int4{}):
			val := srcValue.Interface().(int32)
			dstValue.Set(reflect.ValueOf(toPgInt32(val)))

		case reflect.TypeOf(pgtype.Int8{}):
			val := srcValue.Interface().(int64)
			dstValue.Set(reflect.ValueOf(toPgInt64(val)))

		default:
			if dstValue.Type() == srcValue.Type() {
				dstValue.Set(srcValue)
			}
		}
	}

	return nil
}

func MapStructModelToDomain(src, dst any) error {
	srcVal, err := validateSrc(src)
	if err != nil {
		return err
	}

	dstVal, err := validateDst(dst)
	if err != nil {
		return err
	}

	srcType := srcVal.Type()

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		dstField, ok := dstVal.Type().FieldByName(srcField.Name)
		if !ok {
			continue
		}

		srcValue := srcVal.Field(i)
		dstValue := dstVal.FieldByName(srcField.Name)

		if !dstValue.CanSet() {
			continue
		}

		if srcField.Type.Kind() == reflect.Struct && dstField.Type.Kind() == reflect.Struct {
			err := MapStructModelToDomain(srcValue.Addr().Interface(), dstValue.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		switch srcField.Type {
		case reflect.TypeOf(pgtype.Text{}):
			text := srcValue.Interface().(pgtype.Text)
			dstValue.Set(reflect.ValueOf(text.String))

		case reflect.TypeOf(pgtype.UUID{}):
			id := srcValue.Interface().(pgtype.UUID)
			dstValue.Set(reflect.ValueOf(id.String()))

		case reflect.TypeOf(pgtype.Timestamptz{}):
			t := srcValue.Interface().(pgtype.Timestamptz)
			dstValue.Set(reflect.ValueOf(t.Time))

		case reflect.TypeOf(pgtype.Bool{}):
			b := srcValue.Interface().(pgtype.Bool)
			dstValue.Set(reflect.ValueOf(b.Bool))

		case reflect.TypeOf(pgtype.Int2{}):
			n := srcValue.Interface().(pgtype.Int2)
			dstValue.Set(reflect.ValueOf(n.Int16))

		case reflect.TypeOf(pgtype.Int4{}):
			n := srcValue.Interface().(pgtype.Int4)
			dstValue.Set(reflect.ValueOf(n.Int32))

		case reflect.TypeOf(pgtype.Int8{}):
			n := srcValue.Interface().(pgtype.Int8)
			dstValue.Set(reflect.ValueOf(n.Int64))

		default:
			if srcValue.Type() == dstValue.Type() {
				dstValue.Set(srcValue)
			}
		}
	}

	return nil
}

func MapStruct(src, dst any) error {
	srcVal, err := validateSrc(src)
	if err != nil {
		return err
	}
	srcType := srcVal.Type()

	dstVal, err := validateDst(dst)
	if err != nil {
		return err
	}

	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)

		srcValue := srcVal.Field(i)
		dstValue := dstVal.FieldByName(srcField.Name)

		if dstValue.CanSet() {
			dstValue.Set(srcValue)
		}
	}

	return nil
}

func validateSrc(src any) (reflect.Value, error) {
	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return rv, errors.New("src must be a struct or pointer to struct")
	}
	return rv, nil
}

func validateDst(dst any) (reflect.Value, error) {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return dstVal, fmt.Errorf("dst must be a pointer, got %s", dstVal.Kind())
	}
	if dstVal.Type().Elem().Kind() != reflect.Struct {
		return dstVal, fmt.Errorf("")
	}
	if dstVal.IsNil() {
		return dstVal, fmt.Errorf("dst cannot be nil")
	}

	return dstVal.Elem(), nil
}
