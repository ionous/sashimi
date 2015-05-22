package event

//
type Phase int

const (
	CapturingPhase Phase = iota
	TargetPhase
	BubblingPhase
)

func (this Phase) String() string {
	switch this {
	case CapturingPhase:
		return "CapturingPhase"
	case TargetPhase:
		return "TargetPhase"
	case BubblingPhase:
		return "BubblingPhase"
	}
	return "UnknownPhase"
}
