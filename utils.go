package routerosclient

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-routeros/routeros/proto"
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

func setFieldsFromReply(i interface{}, s *proto.Sentence) (interface{}, error) {
	v := reflect.ValueOf(i).Elem()

	if v.Kind() == reflect.Struct {
		for j := 0; j < v.NumField(); j++ {
			ftag := v.Type().Field(j).Tag.Get("ros")
			fname := v.Type().Field(j).Name
			fval := v.Field(j)

			if s.Map[ftag] != "" {
				if fval.CanSet() && fval.IsValid() {
					switch fval.Kind() {
					case reflect.Bool:
						newVal, err := strconv.ParseBool(s.Map[ftag])
						if err != nil {
							return nil, err
						}
						fval.SetBool(newVal)
					case reflect.Int:
						newVal, err := strconv.ParseInt(s.Map[ftag], 0, 0)
						if !fval.OverflowInt(newVal) {
							fval.SetInt(newVal)
						} else {
							return nil, err
						}
					default:
						fval.SetString(s.Map[ftag])
					}
				} else {
					log.Printf("[W] field `%v` is not settable (ignoring)", fname)
				}
			} else {
				log.Printf("[W] attribute `%v` has no value (ignoring)", ftag)
			}
		}
	}

	return i, nil
}
