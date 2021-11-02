package wikidump

import (
	"fmt"
	"reflect"
	"strconv"
)

func Set(ptr interface{}, tag string) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get(tag); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}
		}
	}
	return nil
}

func setField(field reflect.Value, defaultVal string) error {
	if !field.CanSet() {
		return fmt.Errorf("can't set value")
	}
	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	}
	return nil
}
