package extengine

type AnalysisRequest struct {
	RequestID  string
	EngineID   string
	MultiPV    int
	InitialFEN string
	Moves      []string
	Depth      int
	MoveTime   int
}

type AnalysisResponse struct {
	RequestID string
	Line      string
	End       bool
}
