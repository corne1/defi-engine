package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/corne1/defi-engine/internal/app/dto"
	"github.com/corne1/defi-engine/internal/state"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *dto.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*dto.Transaction, error)
	UpdateState(
		ctx context.Context,
		id uuid.UUID,
		from state.TxState,
		to state.TxState,
	) error
}
