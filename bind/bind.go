package bind

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	// DefaultTag tag
	DefaultTag = "json"
	// IgnoreTag ignore
	IgnoreTag = "-"
)

// Bind parse querystring to interface{}, defaultly json tag.
func Bind(query string, i interface{}) error {
	return BindWithTag(query, i, DefaultTag)
}

// BindWithTag parse query to interface{} by specific tag .
func BindWithTag(query string, i interface{}, tag string) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return errors.New("target must be ptr")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("target must be struct")
	}

	t := v.Type()
	num := t.NumField()
	m := make(map[string]reflect.Value, num)
	parseType(v, m, tag)

	// TODO: error handle
	return parse(m, query)
}

func parseType(v reflect.Value, vals map[string]reflect.Value, tag string) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Tag.Get(tag)
		if name == IgnoreTag {
			continue
		}
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			parseType(v.Field(i), vals, tag)
		}
		if len(name) == 0 {
			name = field.Name
		}
		vals[name] = v.Field(i)
	}
}

func parse(m map[string]reflect.Value, query string) error {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		val := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, val = key[:i], key[i+1:]
		}
		if val == "" {
			continue
		}

		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			// TODO: error handle
			continue
		}
		val, err1 = url.QueryUnescape(val)
		if err1 != nil {
			// TODO: error handle
			continue
		}

		// bind value
		if field, exist := m[key]; exist {
			_ = setValue(field, val)
		}
	}
	return nil
}

func setValue(field reflect.Value, val string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Bool:
		if b, err := strconv.ParseBool(val); err == nil {
			field.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		if i64, err := strconv.ParseInt(val, 10, 0); err == nil {
			field.SetInt(i64)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		if u64, err := strconv.ParseUint(val, 10, 0); err == nil {
			field.SetUint(u64)
		}
	case reflect.Float32, reflect.Float64:
		if f64, err := strconv.ParseFloat(val, 64); err == nil {
			field.SetFloat(f64)
		}
	default:
	}
	return nil
}
