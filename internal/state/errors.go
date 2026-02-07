package state

import "fmt"

type InvalidTransitionError struct {
	From TxState
	To   TxState
}

func (e InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid transition from %s to %s", e.From, e.To)
}
