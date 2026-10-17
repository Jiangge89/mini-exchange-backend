package store

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) GetCheckpoint(ctx context.Context) (*Checkpoint, error) {
	var cp Checkpoint
	var blockHash []byte
	err := p.pool.QueryRow(ctx,
		`SELECT block_number, block_hash FROM checkpoints WHERE id = 1`,
	).Scan(&cp.BlockNumber, &blockHash)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	cp.BlockHash = common.BytesToHash(blockHash)
	return &cp, nil
}

func (p *Postgres) SaveBlock(_ context.Context, _ Checkpoint, _ []Transfer) error {
	return errors.New("not implemented")
}

func (p *Postgres) RollbackToBlock(_ context.Context, _ uint64) error {
	return errors.New("not implemented")
}

func (p *Postgres) GetTransfers(ctx context.Context, address common.Address, limit, offset int) ([]Transfer, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT id, block_number, block_hash, tx_hash, log_index, from_addr, to_addr, value, confirmation_level, created_at
		 FROM transfers
		 WHERE from_addr = $1 OR to_addr = $1
		 ORDER BY block_number DESC, log_index DESC
		 LIMIT $2 OFFSET $3`,
		address.Bytes(), limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var t Transfer
		var blockHash, txHash, from, to []byte
		var value string
		err := rows.Scan(&t.ID, &t.BlockNumber, &blockHash, &txHash, &t.LogIndex, &from, &to, &value, &t.ConfirmationLevel, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		t.BlockHash = common.BytesToHash(blockHash)
		t.TxHash = common.BytesToHash(txHash)
		t.From = common.BytesToAddress(from)
		t.To = common.BytesToAddress(to)
		t.Value, _ = new(big.Int).SetString(value, 10)
		transfers = append(transfers, t)
	}
	return transfers, rows.Err()
}

func (p *Postgres) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	var value string
	err := p.pool.QueryRow(ctx,
		`SELECT amount FROM balances WHERE address = $1`,
		address.Bytes(),
	).Scan(&value)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return big.NewInt(0), nil
		}
		return nil, err
	}
	bal, _ := new(big.Int).SetString(value, 10)
	return bal, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}
