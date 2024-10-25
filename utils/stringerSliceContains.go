package utils

import (
	"fmt"
	"strings"
)

func StringerSliceContains[T fmt.Stringer](s []T, v string) (T, bool) {
	for _, elem := range s {
		if strings.EqualFold(elem.String(), v) {
			return elem, true
		}
	}

	var zero T

	return zero, false
}
