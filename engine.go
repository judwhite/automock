package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"stockhuman/bitboard"
	"stockhuman/chessdb"
	"stockhuman/lichess"
)

type Engine struct {
	UCIOptions []UCIOption

	Hash    int
	Threads int
	MultiPV int

	LichessSpeeds    lichess.Speeds
	LichessRatingMin lichess.Rating
	LichessRatingMax lichess.Rating
	LichessSince     lichess.Date
	LichessUntil     lichess.Date

	fen         string
	moves       []string
	positionMtx sync.RWMutex
}

func NewEngine() *Engine {
	const (
		defaultHash             = 16
		defaultThreads          = 1
		defaultMultiPV          = 1
		defaultLichessSpeeds    = "ultraBullet,bullet,blitz,rapid,classical,correspondence"
		defaultLichessRatingMin = 1600
		defaultLichessRatingMax = 2500
		defaultLichessSince     = "2012-12"
		defaultLichessUntil     = "3000-12"
	)

	e := Engine{
		fen: "startpos",
		UCIOptions: []UCIOption{
			{
				Name:    "Hash",
				Type:    "spin",
				Default: strconv.Itoa(defaultHash),
				Min:     1,
				Max:     33554432,
			},
			{
				Name:    "Threads",
				Type:    "spin",
				Default: strconv.Itoa(defaultThreads),
				Min:     1,
				Max:     1024,
			},
			{
				Name:    "MultiPV",
				Type:    "spin",
				Default: strconv.Itoa(defaultMultiPV),
				Min:     1,
				Max:     256,
			},
			{
				Name:    "Lichess_Speeds",
				Type:    "string",
				Default: defaultLichessSpeeds,
			},
			{
				Name:      "Lichess_Rating_Min",
				Type:      "combo",
				Default:   strconv.Itoa(defaultLichessRatingMin),
				ComboVars: []string{"0", "1000", "1200", "1400", "1600", "1800", "2000", "2200", "2500"},
			},
			{
				Name:      "Lichess_Rating_Max",
				Type:      "combo",
				Default:   strconv.Itoa(defaultLichessRatingMax),
				ComboVars: []string{"0", "1000", "1200", "1400", "1600", "1800", "2000", "2200", "2500"},
			},
			{
				Name:    "Lichess_Since",
				Type:    "string",
				Default: defaultLichessSince,
			},
			{
				Name:    "Lichess_Until",
				Type:    "string",
				Default: defaultLichessUntil,
			},
		},
	}

	for _, uciOption := range e.UCIOptions {
		setOption := fmt.Sprintf("setoption name %s value %s", uciOption.Name, uciOption.Default)
		e.handleSetOption(setOption)
	}

	// runtime self-check. could be moved to a unit test, but, y'know.
	if e.Hash != defaultHash {
		panic(fmt.Errorf("field Hash '%d' != default '%d'", e.Hash, defaultHash))
	}
	if e.Threads != defaultThreads {
		panic(fmt.Errorf("field Threads '%d' != default '%d'", e.Threads, defaultThreads))
	}
	if e.MultiPV != defaultMultiPV {
		panic(fmt.Errorf("field MultiPV '%d' != default '%d'", e.MultiPV, defaultMultiPV))
	}
	if e.LichessSpeeds.String() != defaultLichessSpeeds {
		panic(fmt.Errorf("field LichessSpeeds '%s' != default '%s'", e.LichessSpeeds, defaultLichessSpeeds))
	}
	if e.LichessRatingMin != defaultLichessRatingMin {
		panic(fmt.Errorf("field LichessRatingMin '%d' != default '%d'", e.LichessRatingMin, defaultLichessRatingMin))
	}
	if e.LichessRatingMax != defaultLichessRatingMax {
		panic(fmt.Errorf("field LichessRatingMax '%d' != default '%d'", e.LichessRatingMax, defaultLichessRatingMax))
	}
	if e.LichessSince.String() != defaultLichessSince {
		panic(fmt.Errorf("field LichessSince '%s' != default '%s'", e.LichessSince.String(), defaultLichessSince))
	}
	if e.LichessUntil.String() != defaultLichessUntil {
		panic(fmt.Errorf("field LichessUntil '%s' != default '%s'", e.LichessUntil.String(), defaultLichessUntil))
	}

	return &e
}

func (e *Engine) ParseInput(line string) {
	line = strings.TrimSpace(line)

	parts := strings.Split(line, " ")
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	switch command {
	case "uci":
		e.handleUCI()
	case "setoption":
		e.handleSetOption(line)
	case "ucinewgame":
		e.handleUCINewGame()
	case "isready":
		e.handleIsReady()
	case "position":
		e.handlePosition(line)
	case "go":
		e.handleGo(line)
	case "ponderhit":
		uciWriteLine("info string 'ponderhit' TODO")
	case "stop":
		uciWriteLine("info string 'stop' TODO")
	case "show":
		e.handleShow()
	case "d":
		e.handleD()
	case "quit":
		os.Exit(1)
	default:
		if command != "" {
			uciWriteLine(fmt.Sprintf("info string unknown command '%s'", command))
		}
	}
}

func (e *Engine) handleUCI() {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("id name %s %s\n", EngineName, Version))
	for _, uciOption := range e.UCIOptions {
		sb.WriteString(uciOption.String())
		sb.WriteByte('\n')
	}
	sb.WriteString("uciok")

	uciWriteLine(sb.String())
}

func (e *Engine) handleUCINewGame() {
	e.handlePosition("position startpos")
}

func (e *Engine) handleIsReady() {
	uciWriteLine("readyok")
}

func (e *Engine) handleSetOption(line string) {
	name, value, err := parseSetOption(line)
	if err != nil {
		return
	}

	// check setoption name found in list
	var uciOption UCIOption
	for i := 0; i < len(e.UCIOptions); i++ {
		if strings.EqualFold(e.UCIOptions[i].Name, name) {
			uciOption = e.UCIOptions[i]
			break
		}
	}

	// setoption name not found in list
	if uciOption.Name == "" {
		return
	}

	switch uciOption.Type {
	case "spin":
		// int
		n, err := strconv.Atoi(value)
		if err != nil {
			return
		}

		if n < uciOption.Min || n > uciOption.Max {
			return
		}

		switch strings.ToLower(uciOption.Name) {
		case "hash":
			e.Hash = n
		case "threads":
			e.Threads = n
		case "multipv":
			e.MultiPV = n
		}

	case "string":
		switch strings.ToLower(uciOption.Name) {
		case "lichess_speeds":
			speeds := strings.Split(value, ",")
			lichessSpeeds := make(lichess.Speeds, 0, len(speeds))
			for _, speed := range speeds {
				// check valid enum value
				canonical, ok := lichess.ValidSpeeds.Contains(speed)
				if !ok {
					continue
				}
				// check not duplicate
				if _, ok := lichessSpeeds.Contains(speed); ok {
					continue
				}

				lichessSpeeds = append(lichessSpeeds, canonical)
			}

			if len(lichessSpeeds) == 0 {
				return
			}

			sort.Sort(lichessSpeeds)

			e.LichessSpeeds = lichessSpeeds

		case "lichess_since":
			dt, err := time.Parse("2006-01", value)
			if err != nil {
				return
			}
			e.LichessSince = lichess.Date{Year: dt.Year(), Month: int(dt.Month())}
		case "lichess_until":
			dt, err := time.Parse("2006-01", value)
			if err != nil {
				return
			}
			e.LichessUntil = lichess.Date{Year: dt.Year(), Month: int(dt.Month())}
		}

	case "combo":
		var selectedValue string
		for _, comboValue := range uciOption.ComboVars {
			if strings.EqualFold(value, comboValue) {
				selectedValue = comboValue
				break
			}
		}
		if selectedValue == "" {
			return
		}

		switch strings.ToLower(uciOption.Name) {
		case "lichess_rating_min":
			canonical, ok := lichess.ValidRatings.Contains(selectedValue)
			if !ok {
				return
			}
			e.LichessRatingMin = canonical
		case "lichess_rating_max":
			canonical, ok := lichess.ValidRatings.Contains(selectedValue)
			if !ok {
				return
			}
			e.LichessRatingMax = canonical
		}
	}
}

func parseSetOption(line string) (string, string, error) {
	parts := strings.Split(line, " ")
	// setoption name <name> value <value>
	if len(parts) < 5 {
		return "", "", xerrors.Errorf("setoption string has too few parameters: '%s'", line)
	}
	if parts[0] != "setoption" || parts[1] != "name" {
		return "", "", xerrors.Errorf("setoption string does not start with 'setoption name': '%s'", line)
	}

	optionName := parts[2]
	i := 3
	for ; i < len(parts); i++ {
		if parts[i] == "value" {
			break
		}
		optionName += " " + parts[i]
	}

	i++
	if i >= len(parts) {
		return "", "", xerrors.Errorf("setoption string does not contain 'value': '%s'", line)
	}

	optionValue := strings.Join(parts[i:], " ")

	if optionValue == "<empty>" {
		optionValue = ""
	}

	return optionName, optionValue, nil
}

func (e *Engine) handleShow() {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("info string option name %s value %d\n", "Hash", e.Hash))
	sb.WriteString(fmt.Sprintf("info string option name %s value %d\n", "Threads", e.Threads))
	sb.WriteString(fmt.Sprintf("info string option name %s value %d\n", "MultiPV", e.MultiPV))
	sb.WriteString(fmt.Sprintf("info string option name %s value %s\n", "Lichess_Speeds", e.LichessSpeeds.String()))
	sb.WriteString(fmt.Sprintf("info string option name %s value %d\n", "Lichess_Rating_Min", e.LichessRatingMin))
	sb.WriteString(fmt.Sprintf("info string option name %s value %d\n", "Lichess_Rating_Max", e.LichessRatingMax))
	sb.WriteString(fmt.Sprintf("info string option name %s value %s\n", "Lichess_Since", e.LichessSince.String()))
	sb.WriteString(fmt.Sprintf("info string option name %s value %s\n", "Lichess_Until", e.LichessUntil.String()))

	uciWriteLine(sb.String())
}

func (e *Engine) handlePosition(line string) {
	// position startpos
	// position startpos moves <moves list>
	// position fen <fen>
	// position fen <fen> moves <moves list>

	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return
	}

	var (
		fen   string
		moves []string
	)

	i := 2
	if parts[1] == "startpos" {
		fen = bitboard.StartPos
	} else if parts[1] == "fen" {
		for ; i < len(parts); i++ {
			if parts[i] == "moves" {
				break
			}
		}
		j := i
		if j > 8 {
			j = 8
		}

		fen = strings.Join(parts[2:j], " ")
	} else {
		return
	}

	if i < len(parts)-1 && parts[i] == "moves" {
		moves = parts[i+1:]
	}

	e.positionMtx.Lock()
	e.fen = fen
	e.moves = moves
	e.positionMtx.Unlock()
}

func (e *Engine) readPosition() (string, []string) {
	var fen string
	var moves []string

	e.positionMtx.RLock()
	fen = e.fen
	moves = make([]string, len(e.moves))
	copy(moves, e.moves)
	e.positionMtx.RUnlock()

	return fen, moves
}

func (e *Engine) handleD() {
	fen, moves := e.readPosition()

	var sb strings.Builder
	sb.WriteString("info string position fen ")
	sb.WriteString(fen)
	if len(moves) > 0 {
		sb.WriteString(" moves ")
		sb.WriteString(strings.Join(moves, " "))

		bb, err := bitboard.ParseFEN(fen)
		if err == nil {
			bb, err = bb.Apply(moves)
			if err == nil {
				sb.WriteByte('\n')
				sb.WriteString("info string position fen ")
				sb.WriteString(bb.FEN())
			}
		}
	}

	uciWriteLine(sb.String())
}

func (e *Engine) handleGo(line string) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 1250*time.Millisecond)
	defer cancel()

	args, parseErr := parseGo(line)
	if parseErr != nil {
		uciWriteLine(fmt.Sprintf("info string %s", parseErr.Error()))
		return
	}

	_ = args // TODO: maybe handle the args?

	startFEN, moves := e.readPosition()

	bb, err := bitboard.ParseFEN(startFEN)
	if err != nil {
		uciWriteLine(fmt.Sprintf("info string %s", err.Error()))
		panic(err)
	}
	bb, err = bb.Apply(moves)
	if err != nil {
		uciWriteLine(fmt.Sprintf("info string %s", err.Error()))
		panic(err)
	}

	fen := bb.FEN()

	var wg sync.WaitGroup
	wg.Add(3)

	var (
		suggestedMove lichess.OpeningExplorerMove
		cloudEval     lichess.CloudEvalResponse
		queryAll      chessdb.QueryAllResponse
	)

	go func() {
		defer wg.Done()

		var lichessErr error

		suggestedMove, lichessErr = e.searchLichess(ctx, startFEN, moves)
		if lichessErr != nil {
			// TODO: don't panic
			uciWriteLine(fmt.Sprintf("info string lichess api error: %s", lichessErr.Error()))
			panic(lichessErr)
		}
	}()

	go func() {
		defer wg.Done()

		const multiPV = 3
		var cloudErr error

		cloudEval, cloudErr = lichess.GetCloudEval(ctx, fen, multiPV)
		if cloudErr != nil {
			// TODO: write warning?
			//uciWriteLine(fmt.Sprintf("info string cloudeval api error: %s", cloudErr.Error()))
		}
	}()

	go func() {
		defer wg.Done()

		var chessdbErr error

		queryAll, chessdbErr = chessdb.QueryAll(ctx, fen)
		if chessdbErr != nil {
			// TODO: write warning?
			//uciWriteLine(fmt.Sprintf("info string chessdb api error: %s", chessdbErr.Error()))
		}
	}()

	wg.Wait()

	uci := suggestedMove.UCI

	var cp, mate int

	for _, pv := range cloudEval.PVs {
		pvUCI := strings.Split(pv.MovesUCI, " ")[0]
		if pvUCI == uci {
			cp, mate = pv.CP, pv.Mate
			break
		}
	}

	// allow chessdb to clobber lichess cloud evals
	for _, chessDBMove := range queryAll.Moves {
		if chessDBMove.UCI == suggestedMove.UCI || chessDBMove.SAN == suggestedMove.SAN {
			cp = chessDBMove.Score
			break
		}
	}

	var score string
	if mate == 0 {
		score = fmt.Sprintf("cp %d", cp)
	} else {
		score = fmt.Sprintf("mate %d", mate)
	}

	ms := time.Since(start).Milliseconds()
	msg := fmt.Sprintf("info depth %d time %d score %s pv %s\n"+
		"bestmove %s\n",
		18, ms, score, uci,
		uci,
	)
	uciWriteLine(msg)
}

type GoArgs struct {
	SearchMoves []string
	Ponder      bool
	WTime       int
	BTime       int
	WInc        int
	BInc        int
	MovesToGo   int
	Depth       int
	Nodes       int
	Mate        int
	MoveTime    int
	Infinite    bool
}

func parseGo(line string) (GoArgs, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return GoArgs{}, xerrors.Errorf("go string has too few parameters: '%s'", line)
	}

	var goArgs GoArgs
	i := 1

	getInt := func(cmd string) (int, error) {
		if i >= len(parts)-1 {
			return 0, xerrors.Errorf("go string has missing argument after '%s': '%s'", cmd, line)
		}
		n, err := strconv.Atoi(parts[i+1])
		if err != nil || n < 0 {
			return 0, xerrors.Errorf("go string has invalid int argument after '%s': '%s'", cmd, line)
		}
		i += 2
		return n, nil
	}

loop:
	for i < len(parts) {
		cmd := strings.ToLower(parts[i])
		switch cmd {
		case "searchmoves":
			if i >= len(parts)-1 {
				return GoArgs{}, xerrors.Errorf("go string has missing argument after '%s': '%s'", cmd, line)
			}

			goArgs.SearchMoves = parts[i+1:]
			break loop
		case "ponder":
			goArgs.Ponder = true
			i++
		case "wtime":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.WTime = n
			}
		case "btime":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.BTime = n
			}
		case "winc":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.WInc = n
			}
		case "binc":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.BInc = n
			}
		case "movestogo":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.MovesToGo = n
			}
		case "depth":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.Depth = n
			}
		case "nodes":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.Nodes = n
			}
		case "mate":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.Mate = n
			}
		case "movetime":
			if n, err := getInt(cmd); err != nil {
				return goArgs, err
			} else {
				goArgs.MoveTime = n
			}
		case "infinite":
			goArgs.Infinite = true
			i++
		default:
			return GoArgs{}, xerrors.Errorf("go string has unrecognized command '%s' at position %d: '%s'", cmd, i, line)
		}
	}

	return goArgs, nil
}

func (e *Engine) searchLichess(ctx context.Context, fen string, moves []string) (lichess.OpeningExplorerMove, error) {
	speeds := e.LichessSpeeds
	minRating := e.LichessRatingMin
	maxRating := e.LichessRatingMax
	since := e.LichessSince
	until := e.LichessUntil

	if int(minRating) > int(maxRating) {
		minRating, maxRating = maxRating, minRating
	}

	var ratings lichess.Ratings
	for _, item := range lichess.ValidRatings {
		if item >= minRating && item <= maxRating {
			ratings = append(ratings, item)
		}
	}

	req := lichess.OpeningExplorerRequest{
		FEN:     fen,
		Play:    strings.Join(moves, ","),
		Speeds:  speeds,
		Ratings: ratings,
		Since:   since,
		Until:   until,
	}

	resp, err := lichess.GetLichessGames(ctx, req)
	if err != nil {
		return lichess.OpeningExplorerMove{}, xerrors.Errorf("%w", err)
	}

	suggestedMove := getSuggestedMove(resp)
	return suggestedMove, nil
}
