package api

import (
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

// ParseConfig parse and load given conf to config struct
func ParseConfig(conf Config, config interface{}) {
	for k, v := range conf {
		if k == "plugin" {
			continue
		}
		fieldName := strcase.ToCamel(fmt.Sprintf("%v", k))
		field := reflect.ValueOf(config).Elem().FieldByName(fieldName)
		fieldType := field.Type()
		val := reflect.ValueOf(v)
		field.Set(val.Convert(fieldType))
	}
}
