package pg

import (
	"fmt"
	"reflect"
)

// Struct parsing by 'db' tag
func ParseDbModel(model interface{}) (columns []string, values []interface{}) {
	// Получаем тип структуры
	t := reflect.TypeOf(model)

	// Получаем значение структуры
	v := reflect.ValueOf(model)

	// Проходим по всем полям структуры
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		// Если поле является структурой рекурсивно идем по ней, результат суммируем
		if t.Field(i).Type.Kind() == reflect.Struct {
			subStruct := v.Field(i).Interface()
			c, v := ParseDbModel(subStruct)
			columns = append(columns, c...)
			values = append(values, v...)
		}

		// Получаем тег db каждого поля
		tag := t.Field(i).Tag.Get("db")
		// Если тег db не пустой, то добавляем его в список полей
		if tag != "" {
			columns = append(columns, tag)
			values = append(values, v.Field(i).Interface())
		}
	}

	return columns, values
}

// StructToDBString ...
func StructToDBString(model interface{}) string {
	value := reflect.ValueOf(model)
	typ := reflect.TypeOf(model)

	fields := make([]string, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")

		if tag != "" {
			fieldValue := value.Field(i).Interface()
			// `{"key": "value"}`
			fields = append(fields, fmt.Sprintf("\"%s\":\"%v\"", tag, fieldValue))
		}
	}
	s := "{"
	for i, f := range fields {
		if i == 0 {
			s = s + f
			continue
		}
		s = s + "," + f
	}
	s = s + "}"
	return s
}
