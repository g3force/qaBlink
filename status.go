package main

type QaBlinkStatusCode int

const (
	STABLE   QaBlinkStatusCode = iota
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