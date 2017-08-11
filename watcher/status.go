package watcher

type QaBlinkStatusCode int

const (
	STABLE QaBlinkStatusCode = iota
	UNSTABLE
	FAILED
	UNKNOWN
	DISABLED
)

type QaBlinkState struct {
	StatusCode QaBlinkStatusCode
	Score      uint8
	Pending    bool
}

type QaBlinkJob interface {
	Update()
	State() QaBlinkState
}

func (code QaBlinkStatusCode) String() string {
	switch code {
	case STABLE:
		return "STABLE"
	case UNSTABLE:
		return "UNSTABLE"
	case FAILED:
		return "FAILED"
	case UNKNOWN:
		return "UNKNOWN"
	case DISABLED:
		return "DISABLED"
	}
	panic("Unknown status code")
}
