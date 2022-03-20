package njson

import (
	"reflect"
	"strings"
)

func validTag(tag string) bool {
	return !(tag == "" || tag == "-")
}

func isStructureType(typ string) (ok bool) {
	switch typ {
	case reflect.Slice.String():
		ok = true
	case reflect.Map.String():
		ok = true
	case reflect.Struct.String():
		ok = true
	default:
		ok = false
	}

	if strings.Contains(typ, "[]") {
		ok = true
	}

	return
}
