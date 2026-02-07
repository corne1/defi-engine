package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/corne1/defi-engine/internal/state"
)

type Transaction struct {
	ID          uuid.UUID
	Hash        string
	State       state.TxState
	BlockNumber *int64

	CreatedAt time.Time
	UpdatedAt time.Time
}
