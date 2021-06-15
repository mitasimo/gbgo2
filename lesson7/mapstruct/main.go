package mapstruct

import (
	"errors"
	"reflect"
)

var (
	errDstIsNil          = errors.New("dst is nil")
	errDstMustBeSettable = errors.New("dst must be settable")
	errDstMustBeStruct   = errors.New("dst must a struct")
)

// MapStruct заполняет структуру dest из мапы src
func MapStruct(dst interface{}, src map[string]interface{}) error {
	if dst == nil {
		return errDstIsNil
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return errDstMustBeSettable
	}

	dstSetVal := dstVal.Elem()
	if dstSetVal.Kind() != reflect.Struct {
		return errDstMustBeStruct
	}

	for i := 0; i < dstSetVal.NumField(); i++ {
		field := dstSetVal.Field(i)
		fieldType := dstSetVal.Type().Field(i)

		// пропустить сложные типы
		fieldKind := field.Kind()
		if fieldKind == reflect.Struct ||
			fieldKind == reflect.Array ||
			fieldKind == reflect.Chan ||
			fieldKind == reflect.Func ||
			fieldKind == reflect.Interface ||
			fieldKind == reflect.Map ||
			fieldKind == reflect.Ptr ||
			fieldKind == reflect.Slice {
			//return errors.New("complex type")
			continue
		}

		// получить из мапы значение, соответвующее имени поля структуры
		mval, ok := src[fieldType.Name]
		if !ok {
			//return errors.New("can not find key")
			continue
		}

		// типы значения в мапе и в поле структуры должны совпадать
		mapVal := reflect.ValueOf(mval)
		if field.Kind() != mapVal.Kind() {
			//return errors.New("diff types")
			continue
		}

		// присвоить полю структуры значение из мапы
		field.Set(mapVal)
	}

	return nil
}
