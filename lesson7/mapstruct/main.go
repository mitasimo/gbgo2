package mapstruct

import (
	"errors"
	"fmt"
	"reflect"
)

// MapStruct заполняет структуру dest из мапы src
func MapStruct(dst interface{}, src map[string]interface{}) error {
	if dst == nil {
		return errors.New("dst is nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return errors.New("dst must be settable")
	}

	dstSetVal := dstVal.Elem()
	if dstSetVal.Kind() != reflect.Struct {
		return errors.New("dst must a struct")
	}

	for i := 0; i < dstSetVal.NumField(); i++ {
		typeField := dstSetVal.Type().Field(i)

		fieldKind := typeField.Type.Kind()
		if fieldKind == reflect.Struct ||
			fieldKind == reflect.Array ||
			fieldKind == reflect.Chan ||
			fieldKind == reflect.Func ||
			fieldKind == reflect.Interface ||
			fieldKind == reflect.Map ||
			fieldKind == reflect.Ptr ||
			fieldKind == reflect.Slice {
			// пропустить сложные типы
			continue
		}

		mval, ok := src[typeField.Name]
		if !ok {
			// в мапе отсутствует ключ, соответвующий имени поля структ
			continue
		}

		mapVal := reflect.ValueOf(mval)
		if typeField.Type.Kind() != mapVal.Kind() {
			continue // различаются типы в мапе и структуре
		}

		// присвоить полю структуры значение из мапы

	}

	fmt.Println(dstSetVal)

	return nil
}
