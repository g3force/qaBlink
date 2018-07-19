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
	Id() string
}

func (code QaBlinkStatusCode) String() string {
	switch code {
	case STABLE:
		return "\033[1;32m STABLE \033[0m"
	case UNSTABLE:
		return "\033[1;33mUNSTABLE\033[0m"
	case FAILED:
		return "\033[1;31m FAILED \033[0m"
	case UNKNOWN:
		return "UNKNOWN "
	case DISABLED:
		return "DISABLED"
	}
	panic("Unknown status code")
}
