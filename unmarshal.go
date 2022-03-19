package njson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pingcap/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	N_JSON_TAG = "njson"
)

type NjsonTag struct {
	Fullname   string
	Format     string
	FormatPath string
}

func NewNjsonTag(njsonTag string) (nTag *NjsonTag) {
	nTag = &NjsonTag{}
	if !strings.Contains(njsonTag, ";") {
		nTag.Fullname = njsonTag
		return
	}

	strArr := strings.Split(njsonTag, ";")
	for _, str := range strArr {
		kvArr := strings.Split(str, ":")
		if len(kvArr) != 2 {
			err := errors.Errorf("invalid tag: %s", str)
			panic(err)
		}
		k := strings.Trim(kvArr[0], " ")
		v := strings.Trim(kvArr[1], " ")
		rv := reflect.Indirect(reflect.ValueOf(nTag))
		rt := rv.Type()
		k = strings.ToUpper(k[0:1]) + k[1:]
		_, ok := rt.FieldByName(k)
		if !ok {
			err := errors.Errorf("invalid tag unsupport field: %s", k)
			panic(err)
		}
		rfv := rv.FieldByName(k)
		if !rfv.CanSet() {
			err := errors.Errorf("filed %s unsuport set: in %#v", k, nTag)
			panic(err)
		}
		rfv.Set(reflect.ValueOf(v))
	}
	if nTag.Format != "" && nTag.FormatPath == "" { // 设置了各式化方式，没有设置格式化路径时，使用fullnae作为格式化路径
		nTag.FormatPath = nTag.Fullname
	}

	return
}

var jsonNumberType = reflect.TypeOf(json.Number(""))

// Unmarshal used to unmarshal nested json using "njson" tag
func Unmarshal(data []byte, v interface{}) (err error) {
	if !gjson.ValidBytes(data) {
		return fmt.Errorf("invalid json: %v", string(data))
	}

	// catch code panic and return error message
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("can't unmarshal to invalid type %v", reflect.TypeOf(v))
	}
	elem := rv.Elem()
	typeOfT := elem.Type()
	formatArr := make(map[string]*NjsonTag)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		// Check that the tag is  "njson", and can be set
		if (!validTag(typeOfT.Field(i), N_JSON_TAG)) || !field.CanSet() {
			continue
		}

		// Assume "njson" by default, but change to "json" if tag matches
		tag := typeOfT.Field(i).Tag.Get(N_JSON_TAG)
		njsonTag := NewNjsonTag(tag)
		if njsonTag.Format != "" {
			formatArr[njsonTag.FormatPath] = njsonTag
		}
	}
	//先处理需要format的字段
	for _, njsonTag := range formatArr {
		fragment := gjson.GetBytes(data, njsonTag.FormatPath)
		formatHandler, ok := formatHandlerMap[njsonTag.Format]
		if !ok {
			err := errors.Errorf("format: %s unsupport!!", njsonTag.Format)
			panic(err)
		}
		newJson := fmt.Sprintf(`{"%s":"%s"}`, njsonTag.FormatPath, fragment.String())
		newFragmentStr, err := formatHandler(newJson)
		if err != nil {
			return err
		}
		data, err = sjson.SetBytes(data, njsonTag.FormatPath, newFragmentStr)
		if err != nil {
			return err
		}
	}

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typeOfT.Field(i)

		// Check that the tag is  "njson", and can be set
		if (!validTag(typeOfT.Field(i), N_JSON_TAG)) || !field.CanSet() {
			continue
		}

		// Assume "njson" by default, but change to "json" if tag matches
		tagStr := typeOfT.Field(i).Tag.Get(N_JSON_TAG)
		njsonTag := NewNjsonTag(tagStr)

		// get field value by tag
		result := gjson.GetBytes(data, njsonTag.Fullname)

		// if field type json.Number
		if v != nil && field.Kind() == reflect.String && field.Type() == jsonNumberType {
			SetValue(v, elem, field, fieldType, result.String(), reflect.String)
			continue
		}

		var value interface{}
		if isStructureType(field.Kind().String()) {
			value = parseStructureType(result, field.Type())
		} else {
			// set field value depend on it's data type
			value = parseDataType(result, field.Kind().String())
		}

		if v != nil {
			SetValue(v, elem, field, fieldType, value, reflect.Interface)
		}
	}

	return
}

func SetValue(dst interface{}, elem reflect.Value, field reflect.Value, fieldType reflect.StructField, result interface{}, reflectType reflect.Kind) (err error) {
	vt := reflect.TypeOf(dst)
	fieldName := fieldType.Name
	setFuncName := fmt.Sprintf("%s%s", "Set", fieldName)
	var args []reflect.Value
	var valStr string
	args = append(args, reflect.ValueOf(dst))
	refMethod, refMehtodOk := vt.MethodByName(setFuncName)
	if reflectType == reflect.String {
		valStr, _ = result.(string)
		args = append(args, reflect.ValueOf(valStr))
	} else {
		args = append(args, reflect.ValueOf(result))
	}

	if !refMehtodOk {
		if reflectType == reflect.String {
			field.SetString(valStr)
		} else {
			field.Set(reflect.ValueOf(result))
		}
		return
	}
	err = CallSetMethod(refMethod, args)
	if err != nil {
		return
	}
	return
}

func CallSetMethod(refMethod reflect.Method, args []reflect.Value) (err error) {
	setFuncName := refMethod.Name
	refOutArr := refMethod.Func.Call(args)
	refOutLen := len(refOutArr)
	switch refOutLen {
	case 0:
	case 1:
		refOut := refOutArr[0]
		if err, ok := refOut.Interface().(error); ok {
			if err != nil {
				return err
			} else {
				panic(fmt.Sprintf("%s return value type except error ,got %#v", setFuncName, refOut.Interface()))
			}
		}
	default:
		panic(fmt.Sprintf("%s expect 0 or 1 return value ,got %d", setFuncName, refOutLen))
	}
	return
}

func unmarshalSlice(results []gjson.Result, field reflect.Type) interface{} {
	newSlice := reflect.MakeSlice(field, 0, 0)

	for i := 0; i < len(results); i++ {

		var value interface{}
		if isStructureType(field.Elem().Kind().String()) {
			value = parseStructureType(results[i], field.Elem())
		} else {
			// set field value depend on it's data type
			value = parseDataType(results[i], field.Elem().String())
		}

		if value != nil {
			newSlice = reflect.Append(newSlice, reflect.ValueOf(value))
		}
	}

	return newSlice.Interface()
}

func unmarshalMap(raw string, field reflect.Type) interface{} {
	m := reflect.New(reflect.MapOf(field.Key(), field.Elem())).Interface()

	err := json.Unmarshal([]byte(raw), m)
	if err != nil {
		panic(err)
	}

	return reflect.Indirect(reflect.ValueOf(m)).Interface()
}

func unmarshalStruct(raw string, field reflect.Type) interface{} {
	v := reflect.New(field).Interface()

	err := Unmarshal([]byte(raw), v)
	if err != nil {
		panic(err)
	}

	return reflect.Indirect(reflect.ValueOf(v)).Interface()
}

var formatHandlerMap = map[string]func(input string) (out interface{}, err error){
	"json": func(input string) (out interface{}, err error) {
		var v map[string]string
		err = json.Unmarshal([]byte(input), &v)
		if err != nil {
			return
		}
		var name string
		for name = range v {
			break
		}
		out = new(interface{})
		err = json.Unmarshal([]byte(v[name]), &out)
		if err != nil {
			return
		}
		return
	},
	"comma": func(input string) (out interface{}, err error) {
		out = strings.Split(input, ",")
		return
	},
}
