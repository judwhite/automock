package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Engine struct {
	UCIOptions []UCIOption

	Hash    int
	Threads int
	MultiPV int

	LichessSpeeds    string
	LichessRatingMin int
	LichessRatingMax int
	LichessModes     string
	LichessSince     string
}

func NewEngine() *Engine {
	const (
		defaultHash             = 16
		defaultThreads          = 1
		defaultMultiPV          = 1
		defaultLichessSpeeds    = "ultraBullet,bullet,blitz,rapid,classical,correspondence"
		defaultLichessRatingMin = 1600
		defaultLichessRatingMax = 2500
		defaultLichessModes     = "rated"
		defaultLichessSince     = "2012-12"
	)

	e := Engine{
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
				Name:      "Lichess_Modes",
				Type:      "combo",
				Default:   defaultLichessModes,
				ComboVars: []string{"rated", "casual", "rated,casual"},
			},
			{
				Name:    "Lichess_Since",
				Type:    "string",
				Default: defaultLichessSince,
			},
		},
	}

	for _, uciOption := range e.UCIOptions {
		setOption := fmt.Sprintf("name %s value %s", uciOption.Name, uciOption.Default)
		e.handleSetOption(strings.Split(setOption, " "))
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
	if e.LichessSpeeds != defaultLichessSpeeds {
		panic(fmt.Errorf("field LichessSpeeds '%s' != default '%s'", e.LichessSpeeds, defaultLichessSpeeds))
	}
	if e.LichessRatingMin != defaultLichessRatingMin {
		panic(fmt.Errorf("field LichessRatingMin '%d' != default '%d'", e.LichessRatingMin, defaultLichessRatingMin))
	}
	if e.LichessRatingMax != defaultLichessRatingMax {
		panic(fmt.Errorf("field LichessRatingMax '%d' != default '%d'", e.LichessRatingMax, defaultLichessRatingMax))
	}
	if e.LichessModes != defaultLichessModes {
		panic(fmt.Errorf("field LichessModes '%s' != default '%s'", e.LichessModes, defaultLichessModes))
	}
	if e.LichessSince != defaultLichessSince {
		panic(fmt.Errorf("field LichessSince '%s' != default '%s'", e.LichessSince, defaultLichessSince))
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
		e.handleSetOption(parts[1:])
	case "ucinewgame":
		// no-op
	case "isready":
		uciWriteLine("readyok")
	case "position":
		uciWriteLine("info string 'position' TODO")
	case "go":
		uciWriteLine("info string 'go' TODO")
	case "ponderhit":
		uciWriteLine("info string 'ponderhit' TODO")
	case "stop":
		uciWriteLine("info string 'stop' TODO")
	case "show":
		e.handleShow()
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
	sb.WriteString(fmt.Sprintf("id name author %s\n", Author))
	for _, uciOption := range e.UCIOptions {
		sb.WriteString(uciOption.String())
		sb.WriteByte('\n')
	}
	sb.WriteString("uciok")

	uciWriteLine(sb.String())
}

func (e *Engine) handleSetOption(args []string) {
	name, value := e.parseSetOption(args)
	if name == "" {
		return
	}

	var uciOption UCIOption
	for i := 0; i < len(e.UCIOptions); i++ {
		if strings.EqualFold(e.UCIOptions[i].Name, name) {
			uciOption = e.UCIOptions[i]
			break
		}
	}

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
			// TODO: validate items in list
			e.LichessSpeeds = value
		case "lichess_since":
			dt, err := time.Parse("2006-01", value)
			if err != nil {
				return
			}
			e.LichessSince = fmt.Sprintf("%04d-%02d", dt.Year(), dt.Month())
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
			n, err := strconv.Atoi(selectedValue)
			if err != nil {
				return
			}
			e.LichessRatingMin = n
		case "lichess_rating_max":
			n, err := strconv.Atoi(selectedValue)
			if err != nil {
				return
			}
			e.LichessRatingMax = n
		case "lichess_modes":
			e.LichessModes = selectedValue
		}
	}
}

func (e *Engine) parseSetOption(args []string) (string, string) {
	if len(args) < 3 {
		return "", ""
	}
	if args[0] != "name" {
		return "", ""
	}

	optionName := args[1]
	i := 2
	for ; i < len(args); i++ {
		if args[i] == "value" {
			break
		}
		optionName += " " + args[i]
	}

	i++
	if i >= len(args) {
		return "", ""
	}

	optionValue := strings.Join(args[i:], " ")

	if optionValue == "<empty>" {
		optionValue = ""
	}

	return optionName, optionValue
}

func (e *Engine) handleShow() {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "Hash", e.Hash))
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "Threads", e.Threads))
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "MultiPV", e.MultiPV))
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "Lichess_Speeds", e.LichessSpeeds))
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "Lichess_Rating_Min", e.LichessRatingMin))
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "Lichess_Rating_Max", e.LichessRatingMax))
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "Lichess_Modes", e.LichessModes))
	sb.WriteString(fmt.Sprintf("info string option %s value %v\n", "Lichess_Since", e.LichessSince))

	uciWriteLine(sb.String())
}
