package errutil

import (
	"strconv"
	"strings"
)

type stringer interface {
	String() string
}

// Join expects each part implements Stringer
func New(parts ...interface{}) joined {
	return joined{
		parts: parts,
	}
}

type joined struct {
	parts []interface{}
}

func (j joined) Error() string {
	strs := []string{}
	for i, p := range j.parts {
		var str string
		switch v := p.(type) {
		case bool:
			str = strconv.FormatBool(v)
		case int:
			str = formatInt(int64(v))
		case int32:
			str = formatInt(int64(v))
		case int64:
			str = formatInt(v)
		case uint:
			str = formatUint(uint64(v))
		case uint32:
			str = formatUint(uint64(v))
		case uint64:
			str = formatUint(v)
		case float32:
			str = formatFloat(float64(v))
		case float64:
			str = formatFloat(v)
		case stringer:
			str = v.String()
		case error:
			str = v.Error()
		case string:
			str = v
		default:
			str = "<???>" + strconv.FormatInt(int64(i), 10)
		}
		strs = append(strs, str)
	}
	return strings.Join(strs, " ")
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(float64(v), 'g', -1, 64)
}

func formatInt(v int64) string {
	return strconv.FormatInt(v, 10)
}

func formatUint(v uint64) string {
	return strconv.FormatUint(v, 16)
}
