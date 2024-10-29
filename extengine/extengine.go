package extengine

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/xerrors"

	"stockhuman/utils"
)

const bufferedChannelSize = 4096

type ExternalEngine struct {
	opts Opts

	name              string
	uciVariant        string
	supportedVariants []string

	process       *exec.Cmd
	lastUsedEpoch int64
	isAlive       int64
	stdin         *bufio.Writer
	stdout        *bufio.Scanner

	stopLock     sync.Mutex
	analysisLock sync.Mutex
	requestID    string

	readStream chan string
}

type Opts struct {
	EnginePath string
	SyzygyPath string
	Hash       int
	Threads    int
	MultiPV    int
}

type SetOption struct {
	Name  string
	Value string
}

func New(opts Opts) (*ExternalEngine, error) {
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, opts.EnginePath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, xerrors.Errorf("error creating stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, xerrors.Errorf("error creating stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, xerrors.Errorf("error executing cmd.Start() for '%s': %w", opts.EnginePath, err)
	}

	name := filepath.Base(opts.EnginePath)

	engine := &ExternalEngine{
		opts:              opts,
		name:              name,
		process:           cmd,
		supportedVariants: []string{},
		lastUsedEpoch:     time.Now().Unix(),
		isAlive:           1,
		stdin:             bufio.NewWriter(stdin),
		stdout:            bufio.NewScanner(stdout),
		readStream:        make(chan string, bufferedChannelSize),
	}

	go func() {
		if err := engine.recv(); err != nil {
			utils.Log(fmt.Sprintf("external engine: error reading from engine: %s", err.Error()))
			panic(err)
		}
		utils.Log("external engine: exited")
		atomic.StoreInt64(&engine.isAlive, 0)
	}()

	name, err = engine.uci()
	if err != nil {
		return nil, xerrors.Errorf("error initializing uci: %w", err)
	}
	if name != "" {
		engine.name = name
	}
	utils.Log(fmt.Sprintf("external engine: name: %s", name))

	// Lc0
	// setoption name UCI_Chess960 value true
	// setoption name UCI_ShowWDL value true
	// setoption name UCI_ShowMovesLeft value true
	// setoption name UCI_AnalyseMode value true
	// error Unknown option: UCI_AnalyseMode
	// setoption name MultiPV value 4
	// setoption name WeightsFile value /home/jud/projects/sf/lc0-weights/t2-768x15x24h-swa-5230000.pb.gz
	// setoption name Backend value cuda-auto
	// setoption name NNCacheSize value 5000000
	// setoption name Threads value 12
	// setoption name RamLimitMb value 16384
	// setoption name SyzygyPath value /home/jud/projects/tablebases/cutechess/:/home/jud/projects/tablebases/cutechess/tb6/
	// setoption name Ponder value false
	// setoption name VerboseMoveStats value true
	// setoption name LogLiveStats value true

	setOptions := []SetOption{
		//{"UCI_AnalyseMode", true},
		//{"UCI_Chess960", true},
		//{"UCI_Variant", "chess"},
		{"Hash", strconv.Itoa(engine.opts.Hash)},
		{"Threads", strconv.Itoa(engine.opts.Threads)},
		{"MultiPV", strconv.Itoa(engine.opts.MultiPV)},
		{"SyzygyPath", engine.opts.SyzygyPath},
		{"Ponder", "false"},
	}

	for _, option := range setOptions {
		if err := engine.setOption(option.Name, option.Value); err != nil {
			return nil, xerrors.Errorf("%w", err)
		}
	}

	if err := engine.isReady(); err != nil {
		return nil, xerrors.Errorf("error initializing engine (isready): %w", err)
	}

	if err := engine.send("ucinewgame"); err != nil {
		return nil, xerrors.Errorf("error initializing engine (ucinewgame): %w", err)
	}

	if err := engine.isReady(); err != nil {
		return nil, xerrors.Errorf("error initializing engine (isready): %w", err)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			if err.Error() == "signal: killed" {
				return
			}
			utils.Log(fmt.Sprintf("external engine: abnormal termination: %s", err.Error()))
			panic(err)
		}
	}()

	return engine, nil
}

func (e *ExternalEngine) IdleTime() time.Duration {
	diff := time.Now().Unix() - atomic.LoadInt64(&e.lastUsedEpoch)
	return time.Duration(diff) * time.Second
}

func (e *ExternalEngine) Terminate() error {
	if !atomic.CompareAndSwapInt64(&e.isAlive, 1, 0) {
		return nil
	}

	if err := e.process.Process.Kill(); err != nil {
		return xerrors.Errorf("failed to terminate engine process: %w", err)
	}
	return nil
}

func (e *ExternalEngine) IsAlive() bool {
	return atomic.LoadInt64(&e.isAlive) == 1
}

func (e *ExternalEngine) IsPath(enginePath string) bool {
	return e.opts.EnginePath == enginePath
}

func (e *ExternalEngine) send(command string) error {
	now := time.Now().Unix()
	atomic.StoreInt64(&e.lastUsedEpoch, now)

	_, err := e.stdin.WriteString(command + "\n")
	if err != nil {
		return xerrors.Errorf("error writing to engine, command: '%s', error: %w", command, err)
	}
	if err := e.stdin.Flush(); err != nil {
		return xerrors.Errorf("error flushing command to engine, command: '%s', error: %w", command, err)
	}

	utils.Log(fmt.Sprintf("external engine: > %s", command))

	return nil
}

func (e *ExternalEngine) recv() error {
	for e.stdout.Scan() {
		line := e.stdout.Text()
		if !strings.HasPrefix(line, "info depth") {
			utils.Log(fmt.Sprintf("external engine: < %s", line))
		}
		e.readStream <- line
	}
	if err := e.stdout.Err(); err != nil {
		if errors.Is(err, os.ErrClosed) {
			return nil
		}
		return xerrors.Errorf("stdout: %w", err)
	}
	return nil
}

func (e *ExternalEngine) uci() (string, error) {
	if err := e.send("uci"); err != nil {
		return "", xerrors.Errorf("%w", err)
	}

	var name string

	for line := range e.readStream {
		parts := strings.Split(line, " ")
		command := parts[0]
		if command == "id" {
			if len(parts) >= 3 && parts[1] == "name" && name == "" {
				name = strings.Join(parts[2:], " ")
			}
		} else if command == "option" {
			for i := 1; i < len(parts); i++ {
				if parts[i] == "name" && i+1 < len(parts) {
					name := parts[i+1]
					i++
					if name == "UCI_Variant" && i+1 < len(parts) && parts[i+1] == "var" {
						i++
						e.supportedVariants = append(e.supportedVariants, parts[i])
					}
				}
			}
		} else if command == "uciok" {
			break
		}
	}

	if len(e.supportedVariants) > 0 {
		//
	} else {
		e.supportedVariants = append(e.supportedVariants, "chess")
	}

	return name, nil
}

func (e *ExternalEngine) isReady() error {
	if err := e.send("isready"); err != nil {
		return xerrors.Errorf("%w", err)
	}

	for line := range e.readStream {
		if line == "readyok" {
			break
		}
	}
	return nil
}

func (e *ExternalEngine) setOption(name string, value any) error {
	var cmd string
	switch value := value.(type) {
	case string:
		cmd = fmt.Sprintf("setoption name %s value %s", name, value)
	case int:
		cmd = fmt.Sprintf("setoption name %s value %d", name, value)
	case bool:
		cmd = fmt.Sprintf("setoption name %s value %s", name, strconv.FormatBool(value))
	default:
		return xerrors.Errorf("error: setoption name %s value %v: unsupported option type: %T", name, value, value)
	}

	return e.send(cmd)
}

func (e *ExternalEngine) Analyze(ctx context.Context, job AnalysisRequest, jobStarted chan struct{}) (<-chan string, error) {
	e.analysisLock.Lock()

	if e.requestID != job.RequestID {
		e.requestID = job.RequestID

		if err := e.isReady(); err != nil {
			e.analysisLock.Unlock()
			return nil, xerrors.Errorf("%w", err)
		}
	}

	optionsChanged := false
	if e.opts.MultiPV != job.MultiPV && job.MultiPV >= 1 {
		if err := e.setOption("MultiPV", strconv.Itoa(job.MultiPV)); err != nil {
			e.analysisLock.Unlock()
			return nil, xerrors.Errorf("%w", err)
		}
		e.opts.MultiPV = job.MultiPV
		optionsChanged = true
	}

	if optionsChanged {
		if err := e.isReady(); err != nil {
			e.analysisLock.Unlock()
			return nil, xerrors.Errorf("%w", err)
		}
	}

	var uciPositionCommand string

	if len(job.Moves) == 0 {
		uciPositionCommand = fmt.Sprintf("position fen %s", job.InitialFEN)
	} else {
		uciPositionCommand = fmt.Sprintf("position fen %s moves %s",
			job.InitialFEN,
			strings.Join(job.Moves, " "),
		)
	}

	if err := e.send(uciPositionCommand); err != nil {
		e.analysisLock.Unlock()
		return nil, xerrors.Errorf("%w", err)
	}
	//zap.L().Info("sent to engine", zap.String("job_id", jobID), zap.String("command", uciPositionCommand))
	if err := e.isReady(); err != nil {
		e.analysisLock.Unlock()
		return nil, xerrors.Errorf("%w", err)
	}

	cmd := "go"

	if job.Depth > 0 {
		cmd += " depth " + strconv.Itoa(job.Depth)
	}
	if job.MoveTime > 0 {
		cmd += " movetime " + strconv.Itoa(job.MoveTime)
	}
	if job.Depth == 0 && job.MoveTime == 0 {
		cmd += " infinite"
	}

	if err := e.send(cmd); err != nil {
		e.analysisLock.Unlock()
		return nil, xerrors.Errorf("%w", err)
	}
	jobStarted <- struct{}{}

	stream := make(chan string, bufferedChannelSize)
	go func() {
		defer func() {
			e.analysisLock.Unlock()
		}()

		if err := e.streamAnalysis(ctx, stream); err != nil {
			//zap.L().Warn("error streaming analysis", zap.String("job_id", jobID), zap.Error(err))
		}
	}()

	return stream, nil
}

func (e *ExternalEngine) streamAnalysis(ctx context.Context, analysisStream chan<- string) (returnError error) {
	defer func() {
		if errors.Is(returnError, context.Canceled) {
			//zap.L().Info("streamAnalysis: context canceled", zap.String("job_id", jobID))
			returnError = nil
		}

		close(analysisStream)

		//zap.L().Debug("analysis ended", zap.String("job_id", jobID))
	}()

	for {
		select {
		case <-ctx.Done():
			//zap.L().Debug("streamAnalysis end: context canceled", zap.String("job_id", jobID))
			if err := e.stop(); err != nil {
				returnError = err
				return
			}
			//zap.L().Debug("streamAnalysis end: stop sent", zap.String("job_id", jobID))
			for line := range e.readStream {
				parts := strings.Split(line, " ")
				command := parts[0]
				if command == "bestmove" {
					//zap.L().Debug("streamAnalysis end: bestmove consumed", zap.String("job_id", jobID))
					return
				}
				if command == "info" {
					if !strings.Contains(line, "score") {
						continue
					}

					analysisStream <- line

					//zap.L().Debug("streamAnalysis end: engine output", zap.String("job_id", jobID), zap.String("line", line))
				} else {
					//zap.L().Warn("streamAnalysis end: : unexpected engine command", zap.String("job_id", jobID), zap.String("command", command))
				}
			}
			return
		case line := <-e.readStream:
			parts := strings.Split(line, " ")
			command := parts[0]

			if command == "bestmove" {
				//zap.L().Debug("bestmove received", zap.String("job_id", jobID), zap.String("line", line))

				if err := e.stop(); err != nil {
					returnError = err
					return
				}

				if err := e.isReady(); err != nil {
					returnError = err
					return
				}
				return
			}

			if command == "info" {
				if !strings.Contains(line, "score") {
					continue
				}

				analysisStream <- line

				//zap.L().Debug("engine output", zap.String("job_id", jobID), zap.String("line", line))
			} else {
				//zap.L().Warn("readStream: unexpected engine command", zap.String("job_id", jobID), zap.String("command", command))
			}
		}
	}
}

func (e *ExternalEngine) stop() error {
	if atomic.LoadInt64(&e.isAlive) == 0 {
		return nil
	}

	e.stopLock.Lock()
	defer e.stopLock.Unlock()

	//zap.L().Debug("sending stop", zap.String("job_id", jobID))
	if err := e.send("stop"); err != nil {
		return xerrors.Errorf("%w", err)
	}
	return nil
}

func (e *ExternalEngine) SetOptions(options []SetOption) error {
	for _, o := range options {
		if err := e.setOption(o.Name, o.Value); err != nil {
			return xerrors.Errorf("%w", err)
		}
	}
	return nil
}
