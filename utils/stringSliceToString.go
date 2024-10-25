package utils

import (
	"fmt"
	"strings"
)

func StringerSliceToString[T fmt.Stringer](s []T) string {
	if len(s) == 0 {
		return ""
	}

	v := make([]string, len(s))
	for i, elem := range s {
		v[i] = elem.String()
	}

	return strings.Join(v, ",")
}
