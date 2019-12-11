package openapi

import (
	"net/http"
	"net/url"
	"reflect"
)

const FIELDMASK = "XXX_RestFieldMask"

func fieldMask(in interface{}) []string {
	v := reflect.ValueOf(in)
	if in == nil || (v.Kind() == reflect.Ptr && v.IsNil()) {
		return nil
	}
	f := v.Elem().FieldByName(FIELDMASK)
	if !f.IsValid() {
		return nil
	}
	if f.Kind() == reflect.Slice && f.IsNil() {
		return []string{}
	}
	if f.Kind() == reflect.Slice && f.Type().Elem().Kind() == reflect.String {
		return f.Interface().([]string)
	} else {
		return nil
	}
}

func withFieldMask(req *http.Request, fieldMask []string) {
	// fieldmask没传的时候， 需要传所有的参数。 否则会有不兼容的契约变更
	if len(fieldMask) == 0 {
		return
	}
	value := url.Values{}
	query := req.URL.Query()
	for _, k := range fieldMask {
		if v, ok := query[k]; ok {
			value[k] = v
		}
	}
	req.URL.RawQuery = value.Encode()
}
