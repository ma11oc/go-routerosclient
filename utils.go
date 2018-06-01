package routerosclient

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func buildCommand(command string, proplist *[]string, attrs *map[string]string, isQuery bool) (string, error) {

	cmd := command

	if proplist != nil {
		cmd += fmt.Sprintf(" =.proplist=%v", strings.Join(*proplist, ","))
	}

	if attrs != nil {
		for k, v := range *attrs {
			if v != "" {
				if isQuery {
					cmd += fmt.Sprintf(" ?=%v=%v", k, v)
				} else {
					cmd += fmt.Sprintf(" =%v=%v", k, v)
				}
			}
		}
	}

	return cmd, nil
}

func buildAttrsFromResource(i interface{}) (map[string]string, error) {
	v := reflect.ValueOf(i).Elem()
	attrs := make(map[string]string)

	for j := 0; j < v.NumField(); j++ {
		fieldTag := v.Type().Field(j).Tag.Get("ros")

		if fieldTag != "" {
			attrs[fieldTag] = fmt.Sprintf("%v", v.Field(j))
		}

	}

	return attrs, nil
}

func setFieldsFromMap(r Resource, m map[string]string) (Resource, error) {

	v := reflect.ValueOf(r).Elem()

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			ftag := v.Type().Field(i).Tag.Get("ros")
			fname := v.Type().Field(i).Name
			fval := v.Field(i)

			if m[ftag] != "" {
				if fval.CanSet() && fval.IsValid() {
					switch fval.Kind() {
					case reflect.Bool:
						newVal, err := strconv.ParseBool(m[ftag])
						if err != nil {
							return nil, err
						}
						fval.SetBool(newVal)
					case reflect.Int:
						newVal, err := strconv.ParseInt(m[ftag], 0, 0)
						if !fval.OverflowInt(newVal) {
							fval.SetInt(newVal)
						} else {
							return nil, err
						}
					default:
						fval.SetString(m[ftag])
					}
				} else {
					log.Printf("[W] field `%v` is not settable (ignoring)", fname)
				}
			} else {
				log.Printf("[W] attribute `%v` has no value (ignoring)", ftag)
			}
		}
	}

	return r, nil
}
