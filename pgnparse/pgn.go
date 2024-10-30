package pgnparse

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"automock/bitboard"
)

type GameResult string

const (
	ResultWhiteWins GameResult = "1-0"
	ResultBlackWins GameResult = "0-1"
	ResultDraw      GameResult = "1/2-1/2"
	ResultUnknown   GameResult = "*"
)

type Color int

const (
	PlayerNotFound Color = 0
	WhitePieces    Color = 1
	BlackPieces    Color = -1
)

type PGN struct {
	Games []*Game
}

type Tags []Tag
type Moves []*Move

type Game struct {
	Tags     Tags
	Comment  string
	Moves    Moves
	Result   GameResult
	WhiteElo int
	BlackElo int

	dateTime time.Time
}

func (m Moves) Strings() []string {
	moves := make([]string, 0, len(m))
	for _, move := range m {
		moves = append(moves, move.UCI)
	}
	return moves
}

//func (m Moves) ToMoveWithVariations() *Move {
//	if len(m) == 0 {
//		return nil
//	}
//
//	parent := m[0]
//	if len(parent.Variations) != 0 {
//		xlog.Panic(fmt.Errorf("moves[0] already has variations; expected flat list of single moves"))
//	}
//
//	for i := 1; i < len(m); i++ {
//		varMove := m[i]
//		if len(varMove.Variations) != 0 {
//			xlog.Panic(fmt.Errorf("moves[%d] has variations; expected flat list of single moves", i))
//		}
//
//		parent.Variations = append(parent.Variations, &Variation{Moves: Moves{varMove}})
//	}
//
//	return parent
//}

type MoveNodes []*MoveNode

type MoveNode struct {
	FENKey string
	Ply    int

	SAN string
	UCI string

	Children MoveNodes
}

func (m Moves) ToNodes() MoveNodes {
	if len(m) == 0 {
		return nil
	}

	return m[0].toNodes(m[1:])
}

func (t Tags) Get(name string) string {
	for _, v := range t {
		if v.Name == name {
			return v.Value
		}
	}
	return ""
}

func (g *Game) Equals(g2 *Game) bool {
	if g2 == nil {
		return false
	}

	if len(g.Tags) != len(g2.Tags) {
		return false
	}

	if len(g.Moves) != len(g2.Moves) {
		return false
	}

	if g.Comment != g2.Comment {
		return false
	}

	if g.Result != g2.Result {
		return false
	}

	if g.WhiteElo != g2.WhiteElo || g.BlackElo != g2.BlackElo {
		return false
	}

	for i := 0; i < len(g.Tags); i++ {
		if !g.Tags[i].Equals(g2.Tags[i]) {
			return false
		}
	}

	for i := 0; i < len(g.Moves); i++ {
		if !g.Moves[i].Equals(g2.Moves[i]) {
			return false
		}
	}

	return true
}

func (g *Game) StringNoTags() string {
	var sb strings.Builder

	if g.Comment != "" {
		sb.WriteString(fmt.Sprintf("{ %s } ", g.Comment))
	}

	for i, move := range g.Moves {
		if move.Ply%2 == 1 {
			sb.WriteString(fmt.Sprintf("%d. ", move.FullMoveNumber()))
		} else if i == 0 || len(g.Moves[i-1].Variations) > 0 {
			sb.WriteString(fmt.Sprintf("%d... ", move.FullMoveNumber()))
		}
		sb.WriteString(move.String())
		sb.WriteRune(' ')
	}
	sb.WriteString(string(g.Result))

	return strings.TrimSpace(sb.String())
}

func (g *Game) String() string {
	if g == nil {
		return ""
	}

	var sb strings.Builder

	if len(g.Tags) > 0 {
		for _, tag := range g.Tags {
			sb.WriteString(fmt.Sprintf("[%s \"%s\"]\n", tag.Name, tag.Value))
		}

		sb.WriteByte('\n')
	}

	sb.WriteString(g.StringNoTags())

	return sb.String()
}

func (g *Game) PlayerColor(name string) Color {
	white := g.Tags.Get("White")
	black := g.Tags.Get("Black")
	if strings.EqualFold(white, name) {
		return WhitePieces
	}
	if strings.EqualFold(black, name) {
		return BlackPieces
	}
	return 0
}

func (g *Game) Date() time.Time {
	if !g.dateTime.IsZero() {
		return g.dateTime
	}

	utcDateTag := g.Tags.Get("UTCDate")
	if len(utcDateTag) == 0 {
		return g.dateTime
	}

	utcTimeTag := g.Tags.Get("UTCTime")

	// if date only
	if len(utcTimeTag) == 0 {
		utcDate, err := time.ParseInLocation("2006.01.02", utcDateTag, time.UTC)
		if err != nil {
			return g.dateTime
		}

		g.dateTime = utcDate
		return g.dateTime
	}

	// date and time
	utcDateTimeTag := fmt.Sprintf("%s %s", utcDateTag, utcTimeTag)
	utcDateTime, err := time.ParseInLocation("2006.01.02 15:04:05", utcDateTimeTag, time.UTC)
	if err != nil {
		// try date only
		utcDate, err := time.ParseInLocation("2006.01.02", utcDateTag, time.UTC)
		if err != nil {
			return g.dateTime
		}

		g.dateTime = utcDate
		return g.dateTime
	}

	g.dateTime = utcDateTime
	return g.dateTime
}

type Tag struct {
	Name  string
	Value string
}

func (t Tag) Equals(t2 Tag) bool {
	return t.Name == t2.Name && t.Value == t2.Value
}

type Move struct {
	FENKey string
	Ply    int

	SAN        string
	UCI        string
	NAGs       []string
	Comment    string
	Variations []*Variation
}

func (m *Move) toNodes(descendants Moves) MoveNodes {
	var nodes MoveNodes

	topNode := &MoveNode{
		FENKey: m.FENKey,
		Ply:    m.Ply,
		SAN:    m.SAN,
		UCI:    m.UCI,
	}

	if len(descendants) > 0 {
		topNode.Children = descendants[0].toNodes(descendants[1:])
	}

	nodes = append(nodes, topNode)

	for _, variation := range m.Variations {
		// variations Moves[0] will never have variations of their own
		varNodes := variation.Moves.ToNodes()
		if len(varNodes) != 1 {
			panic(fmt.Errorf("len(varNodes) = %d, expected 1", len(varNodes)))
		}
		nodes = append(nodes, varNodes[0])
	}

	return nodes
}

func (m *Move) FullMoveNumber() int {
	return (m.Ply + 1) / 2
}

func (m *Move) Equals(m2 *Move) bool {
	if m2 == nil {
		return false
	}

	if m.Ply != m2.Ply ||
		m.SAN != m2.SAN ||
		m.Comment != m2.Comment ||
		!reflect.DeepEqual(m.NAGs, m2.NAGs) {
		return false
	}

	if len(m.Variations) != len(m2.Variations) {
		return false
	}

	for i := 0; i < len(m.Variations); i++ {
		v1 := m.Variations[i]
		v2 := m.Variations[i]

		if len(v1.Comments) != len(v2.Comments) {
			return false
		}

		for j := 0; j < len(v1.Comments); j++ {
			if v1.Comments[j] != v2.Comments[j] {
				return false
			}
		}

		if len(v1.Moves) != len(v2.Moves) {
			return false
		}

		for j := 0; j < len(v1.Moves); j++ {
			m1 := v1.Moves[j]
			m2 := v2.Moves[j]

			if !m1.Equals(m2) {
				return false
			}
		}
	}

	return true
}

func (m *Move) String() string {
	var sb strings.Builder
	m.writeMove(&sb, 1)
	return strings.TrimSpace(sb.String())
}

func writeMove(sb *strings.Builder, san string, nags []string, comment string) {
	sb.WriteString(fmt.Sprintf("%s ", san))
	for _, nag := range nags {
		sb.WriteString(fmt.Sprintf("%s ", nag))
	}
	//if comment != "" {
	//	sb.WriteString(fmt.Sprintf("{ %s } ", comment))
	//}
}

func (m *Move) writeMove(sb *strings.Builder, indent int) bool {
	writeMove(sb, m.SAN, m.NAGs, m.Comment)

	if len(m.Variations) == 0 {
		return false
	}

	for i := 0; i < len(m.Variations); i++ {
		v := m.Variations[i]

		sb.WriteString(fmt.Sprintf("\n%s( ", strings.Repeat(" ", indent*2)))
		var lastHadVariations bool
		for j, move := range v.Moves {
			//if j == 0 {
			//for _, comment := range v.Comments {
			//	sb.WriteString(fmt.Sprintf("{%s} ", comment))
			//}

			if j == 0 || lastHadVariations {
				// TODO: "..." for Black
				if bitboard.PlyToColor(move.Ply) == bitboard.White {
					sb.WriteString(fmt.Sprintf("%d. ", move.FullMoveNumber()))
				} else {
					sb.WriteString(fmt.Sprintf("%d... ", move.FullMoveNumber()))
				}
			}

			lastHadVariations = move.writeMove(sb, indent+1)
			//writeMove(sb, move.SAN, move.NAGs, move.Comment)
		}
		sb.WriteString(")\n")
	}

	return true
}

type Variation struct {
	Comments []string
	Moves    Moves
}
