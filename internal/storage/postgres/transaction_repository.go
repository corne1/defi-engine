package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/corne1/defi-engine/internal/app/dto"
	"github.com/corne1/defi-engine/internal/state"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *dto.Transaction) error {
	query := `
		INSERT INTO transactions (id, hash, state, block_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()

	_, err := r.db.Exec(
		ctx,
		query,
		tx.ID,
		tx.Hash,
		tx.State,
		tx.BlockNumber,
		now,
		now,
	)

	return err
}

func (r *TransactionRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*dto.Transaction, error) {
	query := `
		SELECT id, hash, state, block_number, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var tx dto.Transaction

	err := r.db.QueryRow(ctx, query, id).
		Scan(
			&tx.ID,
			&tx.Hash,
			&tx.State,
			&tx.BlockNumber,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &tx, err
}

func (r *TransactionRepository) UpdateState(
	ctx context.Context,
	id uuid.UUID,
	from state.TxState,
	to state.TxState,
) error {
	if !state.CanTransition(from, to) {
		return state.InvalidTransitionError{From: from, To: to}
	}

	query := `
		UPDATE transactions
		SET state = $1, updated_at = NOW()
		WHERE id = $2 AND state = $3
	`

	cmd, err := r.db.Exec(ctx, query, to, id, from)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("state was changed concurrently")
	}

	return nil
}
