// Package queryables is a helper package which collect HTTP query
//  and possibly transform them into actual database query
package queryables

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Info of query
type Info struct {
	_key string

	DtoKey    string       // DtoKey is the field name which appears on JSON / Publicly known
	DaoKey    string       // DaoKey is the field name which is stored in database
	TypeOf    reflect.Kind // TypeOf data
	Transform Transform    // Transform the query into actual database query
	Default   string       // Default value if query does not exists in the request
}

// Key returns either DtoKey if exists, or DaoKey if DtoKey doesn't exists
func (i *Info) Key() string {
	// if cache is not empty
	if i._key != "" {
		return i._key
	}

	// store key to cache
	i._key = i.DtoKey
	if i._key == "" {
		i._key = i.DaoKey
	}

	return i._key
}

// Value get query value from HTTP request
// returns nil if both default value and val is empty
// returns actual value parsed into what its' type of
func (i *Info) Value(r *http.Request) interface{} {
	key := i.Key()
	str := r.URL.Query().Get(key)

	// use default if query value doesn't exists
	if str == "" {
		str = i.Default
	}

	// returns nil if both default value and val is empty
	if str == "" {
		return nil
	}

	// convert query value from string to actual value
	var val interface{}
	var err error
	switch i.TypeOf {
	case reflect.Bool:
		val, err = strconv.ParseBool(str)
	case reflect.Float32, reflect.Float64:
		val, err = strconv.ParseFloat(str, 64)
	case reflect.Int, reflect.Int32, reflect.Int64:
		val, err = strconv.Atoi(str)
	case reflect.Array, reflect.Slice:
		val = strings.Split(str, ",")
	default:
		return str
	}

	// returns error if conversion failed
	if err != nil {
		return nil
	}

	// return actual value
	return val
}

//Transform query into another value
type Transform func(key string, value interface{}) (string, interface{})
