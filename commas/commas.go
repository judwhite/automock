package commas

import "fmt"

func Int(v int) string {
	s := fmt.Sprintf("%d", v)
	return String(s)
}

func String(s string) string {
	addNegative := false
	if s[0] == '-' {
		addNegative = true
		s = s[1:]
	}

	pos := len(s) - 3
	for pos > 0 {
		s = s[:pos] + "," + s[pos:]
		pos -= 3
	}

	if addNegative {
		return "-" + s
	}
	return s
}
