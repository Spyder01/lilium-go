package config

import (
	"reflect"

	"github.com/spyder01/lilium-go/pkg/utils/env"
)

func expandEnvFields(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanSet() {
				expandEnvFields(field)
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			expandEnvFields(v.Index(i))
		}

	case reflect.String:
		// Apply your function to string fields
		str := v.String()
		v.SetString(env.ExpandEnvWithDefault(str))
	}
}

func ResolveEnv(cfg *LiliumConfig) {
	expandEnvFields(reflect.ValueOf(cfg))
}
