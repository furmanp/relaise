package internal

import (
	"reflect"
)

func MapConfigToPrompt(cfg Config) NotesPrompt {
	np := NotesPrompt{}
	cfgVal := reflect.ValueOf(cfg)
	cfgType := reflect.TypeOf(cfg)
	npVal := reflect.ValueOf(&np).Elem()

	for i := 0; i < cfgVal.NumField(); i++ {
		field := cfgType.Field(i)
		fieldName := field.Name
		npField := npVal.FieldByName(fieldName)

		if npField.IsValid() && npField.CanSet() && npField.Type() == cfgVal.Field(i).Type() {
			npField.Set(cfgVal.Field(i))
		}
	}

	return np
}
