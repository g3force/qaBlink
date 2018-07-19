package watcher

import "fmt"

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
		return "\033[1;31mUNKNOWN \033[0m"
	case DISABLED:
		return "\033[1;33mDISABLED\033[0m"
	}
	panic("Unknown status code")
}

func (s QaBlinkState) String() string {
	if s.Pending {
		return fmt.Sprintf("\033[1;46m%v\033[0m", s.StatusCode)
	}
	return s.StatusCode.String()
}
