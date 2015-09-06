// Package recode rewrites json-interface objects ( with golang upper case keys )
// to Json dictionaries ( and lowercase keys. )
package recode

import "strings"

type jsonSlice []interface{}
type jsonMap map[string]interface{}

func Dict(src jsonMap) jsonMap {
	dst := map[string]interface{}{}
	for k, v := range src {
		lowerKey := strings.ToLower(k[0:1]) + k[1:]
		dst[lowerKey] = Value(v)
	}
	return dst
}

func Array(src jsonSlice) jsonSlice {
	for i, v := range src {
		src[i] = Value(v)
	}
	return src
}

func Value(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return Array(v)
	case map[string]interface{}:
		return Dict(v)
	}
	return v
}
