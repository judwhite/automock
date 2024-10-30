package pgnparse

import "fmt"

type item struct {
	typ itemType
	val []byte
	pos int
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return string(i.val)
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.40q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}
