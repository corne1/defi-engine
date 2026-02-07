package state

type TxState string

const (
	TxStatePending   TxState = "pending"   // создана, но не отправлена
	TxStateSent      TxState = "sent"      // отправлена в сеть
	TxStateConfirmed TxState = "confirmed" // подтверждена
	TxStateFailed    TxState = "failed"    // ошибка (revert, gas, etc)
	TxStateReverted  TxState = "reverted"  // откат из-за reorg
)
