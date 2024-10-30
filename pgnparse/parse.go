package pgnparse

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/xerrors"

	"automock/bitboard"
)

func ParseReader(r io.Reader) (*PGN, error) {
	const bufSize = 16384

	newGameMarker := []byte{'\n', '\n', '['}

	var buf [bufSize]byte
	bigBuf := make([]byte, 0, bufSize*2)

	var pgn PGN

	var eof bool
	for !eof {
		n, err := r.Read(buf[:])
		if err != nil && err != io.EOF {
			return &pgn, err
		}
		if err == io.EOF {
			eof = true
		}

		bigBuf = append(bigBuf, buf[:n]...)
		idx := bytes.LastIndex(bigBuf, newGameMarker)
		if idx == -1 {
			if !eof {
				continue
			} else {
				idx = len(bigBuf)
			}
		}

		input := bigBuf[:idx]

		result, err := parse(input)
		if err != nil {
			return &pgn, err
		}
		pgn.Games = append(pgn.Games, result.Games...)

		if !eof {
			start := idx + len(newGameMarker) - 1
			copy(bigBuf, bigBuf[start:])
			bigBuf = bigBuf[:len(bigBuf)-start]
		} else {
			bigBuf = bigBuf[:0]
		}
	}

	if len(bigBuf) > 0 {
		return &pgn, fmt.Errorf("unprocessed input: '''%s'''", string(bigBuf))
	}

	return &pgn, nil
}

func (pgn *PGN) HydrateMoves() error {
	var wg sync.WaitGroup
	wg.Add(len(pgn.Games))
	for i := 0; i < len(pgn.Games); i++ {
		go func(game *Game) {
			defer wg.Done()

			var startPos = bitboard.StartPosKey
			if pos := game.Tags.Get("FEN"); pos != "" {
				startPos = pos
			}

			if err := fillMovesUCIs(game.Moves, startPos, 0); err != nil {
				panic(fmt.Errorf("startpos: '%s' %v\ngame:\n%s", startPos, err, game.String()))
			}
		}(pgn.Games[i])
	}

	wg.Wait()

	return nil
}

func Parse(input string) (*PGN, error) {
	pgn, err := parse([]byte(input))
	if err != nil {
		return pgn, err
	}

	for _, game := range pgn.Games {
		var startPos = bitboard.StartPosKey
		if pos := game.Tags.Get("FEN"); pos != "" {
			startPos = pos
		}

		if err := fillMovesUCIs(game.Moves, startPos, 0); err != nil {
			return pgn, fmt.Errorf("startpos: '%s' %v\ngame:\n%s", startPos, err, game.String())
		}
	}

	return pgn, nil
}

func movesToString(moves []*Move) string {
	var sb strings.Builder
	for _, move := range moves {
		if sb.Len() != 0 {
			sb.WriteByte(' ')
		}
		if move.Ply%2 == 1 {
			sb.WriteString(fmt.Sprintf("%d. ", move.FullMoveNumber()))
		}
		sb.WriteString(move.SAN)
	}
	return sb.String()
}

func fillMovesUCIs(moves []*Move, pos string, depth int) error {
	b, err := bitboard.ParseFEN(pos)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	for i, move := range moves {
		move.FENKey = b.FENKey()

		uci, err := b.UCI(move.SAN)
		if err != nil {
			return fmt.Errorf("%s: %v", movesToString(moves[:i+1]), err)
		}
		move.UCI = uci

		if b, err = b.Apply(move.UCI); err != nil {
			return fmt.Errorf("%s: %v", movesToString(moves[:i+1]), err)
		}

		// Variations are children of the current move's parent
		for _, v := range move.Variations {
			if err := fillMovesUCIs(v.Moves, move.FENKey, depth+1); err != nil {
				return fmt.Errorf("%s %v", movesToString(moves[:i+1]), err)
			}
		}
	}
	//fmt.Println()

	return nil
}

func parse(input []byte) (*PGN, error) {
	l, err := lex(input)
	if err != nil {
		return nil, err
	}

	var games []*Game
	items := l.items

	for len(items) > 0 {
		if items[0].typ == itemEOF {
			break
		}

		game, rest, err := parseGame(items)
		if err != nil {
			if game == nil {
				return nil, fmt.Errorf("game #%d: %v", len(games)+1, err)
			}
			return nil, fmt.Errorf("game #%d: tags: %v: %v", len(games)+1, game.Tags, err)
		}

		variant := game.Tags.Get("Variant")
		if variant == "" || strings.EqualFold(variant, "Standard") {
			games = append(games, game)
		}

		items = rest
	}

	return &PGN{Games: games}, nil
}

func parseGame(items []item) (*Game, []item, error) {
	var game Game

	// Parse Tags section
tagLoop:
	for i := 0; i < len(items); i++ {
		item := items[i]

		switch item.typ {

		case itemTagName:
			tagName := item.val
			i++
			item := items[i]
			if item.typ != itemString {
				panic(fmt.Errorf("expected itemString after itemTagName; got: type: %s val: '%s'", item.typ, item.val))
			}
			game.Tags = append(game.Tags, Tag{Name: string(tagName), Value: string(item.val)})

		case itemComment, itemMoveNumber:
			items = items[i:]
			break tagLoop

		case itemGameTermination:
			// aborted game
			game.Result = GameResult(item.val)
			items = items[i+1:]
			return &game, items, nil

		default:
			return &game, items, fmt.Errorf("parseGame (tags/moves): unhandled token. type: %s value: '%s'", item.typ, item.val)
		}
	}

	startFEN := game.Tags.Get("FEN")
	startPly := 1
	if startFEN != "" {
		bb, err := bitboard.ParseFEN(startFEN)
		if err != nil {
			return nil, nil, xerrors.Errorf("error parsing PGN Tag 'FEN': '%s': %w", startFEN, err)
		}
		startPly = bb.FullMoveNumber - (1 - int(bb.ActiveColor))
	}

preGameCommentsLoop:
	for i := 0; i < len(items); i++ {
		item := items[i]

		switch item.typ {
		case itemComment:
			if game.Comment != "" {
				return nil, nil, fmt.Errorf("game comment already set")
			}
			game.Comment = string(item.val)
		case itemMoveNumber:
			items = items[i:]
			break preGameCommentsLoop
		default:
			panic(fmt.Errorf("parseGame (pre-moves): unhandled token. type: %s value: '%s'", item.typ, item.val))
		}
	}

	moves, items, err := parseMoves(items, startPly, 256)
	if err != nil {
		return nil, nil, err
	}
	game.Moves = moves

postGameTerminationLoop:
	for i := 0; i < len(items); i++ {
		item := items[i]

		switch item.typ {
		case itemGameTermination:
			game.Result = GameResult(item.val)
			items = items[i+1:]
			break postGameTerminationLoop
		default:
			panic(fmt.Errorf("parseGame (post-moves): unhandled token. type: %s value: '%s'", item.typ, item.val))
		}
	}

	if game.Result == "" {
		game.Result = ResultUnknown
	}

	if whiteElo := game.Tags.Get("WhiteElo"); whiteElo != "" && whiteElo != "?" {
		if n, err := strconv.Atoi(whiteElo); err == nil {
			game.WhiteElo = n
		}
	}

	if blackElo := game.Tags.Get("BlackElo"); blackElo != "" && blackElo != "?" {
		if n, err := strconv.Atoi(blackElo); err == nil {
			game.BlackElo = n
		}
	}

	return &game, items, nil
}

func parseVariation(items []item, startPly int) (*Variation, []item, error) {
	var variation Variation
	var i int

	for i = 0; i < len(items); i++ {
		item := items[i]

		switch item.typ {
		case itemMoveNumber:
			if len(variation.Moves) != 0 {
				panic(fmt.Errorf("variation already has moves"))
			}

			varMoves, rest, err := parseMoves(items[i:], startPly, 64)
			if err != nil {
				return nil, nil, err
			}

			variation.Moves = varMoves

			i = -1
			items = rest

		case itemComment:
			variation.Comments = append(variation.Comments, string(item.val))

		case itemRightParen:
			return &variation, items[i+1:], nil

		default:
			panic(fmt.Errorf("parseVariation: unhandled token. type: %s value: '%s'", item.typ, item.val))
		}
	}

	panic("unexpected eof")
}

func parseMoves(items []item, startPly, estimatedMoves int) ([]*Move, []item, error) {
	moves := make([]*Move, 0, estimatedMoves)

	var i int
movesLoop:
	for i = 0; i < len(items); i++ {
		item := items[i]

		switch item.typ {

		case itemMoveNumber:
			// safe to ignore

		case itemMoveSAN:
			san := string(item.val)

			// translate gxf8Q+ to gxf8=Q+
			// TODO: slow
			//idx := len(item.val) - 1
			//suffix := item.val[idx]
			//if suffix == '+' || suffix == '#' {
			//	suffix = item.val[idx-1]
			//	idx--
			//}
			//if suffix == 'Q' || suffix == 'R' || suffix == 'B' || suffix == 'N' {
			//	if item.val[idx-1] != '=' {
			//		san = fmt.Sprintf("%s=%s", item.val[:idx], item.val[idx:])
			//	}
			//}

			move := Move{SAN: san}
			move.Ply = startPly + len(moves)
			moves = append(moves, &move)

		case itemMoveNAG:
			if len(moves) == 0 {
				break
			}
			lastMove := moves[len(moves)-1]
			lastMove.NAGs = append(lastMove.NAGs, string(item.val))

		case itemComment:
			if len(moves) == 0 {
				break
			}
			lastMove := moves[len(moves)-1]
			lastMove.Comment = string(item.val)

		case itemGameTermination:
			return moves, items[i:], nil

		case itemLeftParen:
			varStartPly := startPly + len(moves) - 1
			variation, rest, err := parseVariation(items[i+1:], varStartPly)
			if err != nil {
				return nil, nil, err
			}

			lastMove := moves[len(moves)-1]
			lastMove.Variations = append(lastMove.Variations, variation)

			items = rest
			i = -1
			continue

		case itemRightParen:
			return moves, items[i:], nil

		case itemEOF:
			break movesLoop

		default:
			panic(fmt.Errorf("parseMoves: unhandled token. type: %s value: '%s'", item.typ, item.val))
		}
	}

	if i >= len(items) || items[i].typ == itemEOF {
		return moves, nil, nil
	}

	return moves, items[i+1:], nil
}
