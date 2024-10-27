package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"

	"stockhuman/lichess"
)

const (
	EngineName = "Stockhuman"
	Version    = "1.0"
	Author     = "bonerpull"
)

func main() {
	/*
			TODO: 1. Create REPL for stdin/stdout to accept UCI commands and write back status updates
			TODO: 2. Process UCI command: see file uci_to_implement_for_human_engine.txt
			TODO: 3. Store Lichess API results in a DB with a timestamp (to detect stale results), and query parameters.
			TODO:    This will prevent us from querying Lichess too frequently for the same data.
			TODO: 4. On "go <args>" UCI command, we should be able to query the database and Lichess to get a list of MultiPV moves.
			TODO:	 We might need to add a UCI option (setoption <name> value <value>) about the "go <args>" behavior.
			TODO: 	 Originally, it would do the math and return a move according to the distribution; however, it also seems useful
		    TODO:    to return a list of moves based on popularity when MultiPV > 1.
	*/

	//fen := "r1bqkb1r/ppp2ppp/2n2n2/1B1pp3/4P3/P1N2N2/1PPP1PPP/R1BQK2R b KQkq - 1 5" // Gunsberg
	//fen := "r1bqkb1r/ppp2ppp/2n2n2/1B2N3/4p3/P1N5/1PPP1PPP/R1BQK2R b KQkq - 0 6" // Gunsberg
	//fen := "rnbqkb1r/1p2pppp/p2p1n2/8/3NP3/2N5/PPP2PPP/R1BQKB1R w KQkq - 0 6" // Najdorf

	uciWriteLine(fmt.Sprintf("%s %s by %s", EngineName, Version, Author))
	uciLoop()
}

var stdoutMutex sync.Mutex

func uciWriteLine(line string) {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}

	stdoutMutex.Lock()
	fmt.Fprintf(os.Stdout, "%s", line)
	stdoutMutex.Unlock()
}

func uciLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine := NewEngine()

	c := make(chan string, 512)

	go func() {
		defer close(c)
		r := bufio.NewScanner(os.Stdin)

		for r.Scan() {
			select {
			case c <- r.Text():
			case <-ctx.Done():
				return
			}
		}

		if err := r.Err(); err != nil {
			msg := fmt.Sprintf("info ERR: %v", err)
			uciWriteLine(msg)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for line := range c {
			engine.ParseInput(line)
		}
	}()

	wg.Wait()
}

func getSuggestedMove(resp lichess.OpeningExplorerResponse) string {
	var sumMoveTotals int

	type Range struct {
		Lower int
		Upper int
	}

	m := make(map[string]Range)
	sanToUCI := make(map[string]string)

	for _, move := range resp.Moves {
		r := Range{Lower: sumMoveTotals}

		moveTotal := move.Total()
		sumMoveTotals += moveTotal

		r.Upper = sumMoveTotals

		m[move.SAN] = r
		sanToUCI[move.SAN] = move.UCI

		//fmt.Printf("Move: %-7s Games: %11s Popularity: %5.1f%% White: %5.1f%% Draw: %5.1f%% Black: %5.1f%%\n",
		//	move.SAN,
		//	commas.Int(moveTotal),
		//	float64(moveTotal)/float64(sumMoveTotals)*100,
		//	float64(move.White)/float64(moveTotal)*100,
		//	float64(move.Draws)/float64(moveTotal)*100,
		//	float64(move.Black)/float64(moveTotal)*100,
		//)
	}

	// TODO
	if sumMoveTotals == 0 {
		return "0000"
	}

	getRandomMove := func() string {
		n := rand.Intn(sumMoveTotals)
		for k, v := range m {
			if n >= v.Lower && n < v.Upper {
				return k
			}
		}
		panic(fmt.Errorf("couldn't find entry that satisfied n=%d sumMoveTotals=%d", n, sumMoveTotals))
	}

	//const totalSimulationRuns = 1_000_000
	//simulationResults := make(map[string]int)
	//for i := 0; i < totalSimulationRuns; i++ {
	//	san := getRandomMove()
	//	simulationResults[san] += 1
	//}
	//
	//for _, move := range resp.Moves {
	//	r := m[move.SAN]
	//	moveTotal := r.Upper - r.Lower
	//
	//	actualPopularity := float64(moveTotal) / float64(sumMoveTotals) * 100
	//	simulatedPopularity := float64(simulationResults[move.SAN]) / totalSimulationRuns * 100
	//
	//	fmt.Printf("Move: %-7s Act. Popularity: %5.2f%% Sim. Popularity: %5.2f%% Delta: %5.2f%%\n",
	//		move.SAN,
	//		actualPopularity,
	//		simulatedPopularity,
	//		math.Abs(actualPopularity-simulatedPopularity),
	//	)
	//}

	san := getRandomMove()
	return sanToUCI[san]
}
