package pgnparse

import "fmt"

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred: value is text of error
	itemEOF

	itemTagName // tag name in [Name "Value"]
	itemString  // quoted string (excluding quotes)

	itemMoveNumber      // move number
	itemMoveSAN         // move SAN
	itemMoveNAG         // move NAG (Numeric Annotation Glyph)
	itemGameTermination // game termination marker

	itemComment // comment string

	itemLeftParen  // left parentheses: (
	itemRightParen // right parentheses: )
)

func (i itemType) String() string {
	switch i {
	case itemError:
		return "ERROR"
	case itemEOF:
		return "EOF"
	case itemMoveNumber:
		return "MoveNumber"
	case itemMoveSAN:
		return "MoveSAN"
	case itemMoveNAG:
		return "MoveNAG"
	case itemGameTermination:
		return "GameTermination"
	case itemString:
		return "String"
	case itemLeftParen:
		return "LeftParen"
	case itemRightParen:
		return "RightParen"
	case itemComment:
		return "Comment"
	case itemTagName:
		return "TagName"
	}
	return fmt.Sprintf("UNKNOWN - %d", i)
}
