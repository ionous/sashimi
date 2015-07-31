package event

//
type Phase int

const (
	CapturingPhase Phase = iota
	TargetPhase
	BubblingPhase
)

func (phase Phase) String() string {
	switch phase {
	case CapturingPhase:
		return "CapturingPhase"
	case TargetPhase:
		return "TargetPhase"
	case BubblingPhase:
		return "BubblingPhase"
	}
	return "UnknownPhase"
}
