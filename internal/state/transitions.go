package state

var transitions = map[TxState][]TxState{
	TxStatePending: {
		TxStateSent,
		TxStateFailed,
	},
	TxStateSent: {
		TxStateConfirmed,
		TxStateFailed,
		TxStateReverted,
	},
	TxStateConfirmed: {
		TxStateReverted,
	},
	TxStateFailed: {},
	TxStateReverted: {
		TxStateSent, // повторная отправка
	},
}

func CanTransition(from, to TxState) bool {
	allowed, ok := transitions[from]
	if !ok {
		return false
	}

	for _, state := range allowed {
		if state == to {
			return true
		}
	}

	return false
}
