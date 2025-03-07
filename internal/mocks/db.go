package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX is a mock implementation of the database.DBTX interface
type DBTX struct {
	Balance int32
	Err     error
}

func (m *DBTX) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (m *DBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *DBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return &MockRow{Balance: m.Balance, Err: m.Err}
}

// MockRow is a mock implementation of pgx.Row
type MockRow struct {
	Balance int32
	Err     error
}

func (r *MockRow) Scan(dest ...any) error {
	if r.Err != nil {
		return r.Err
	}
	*dest[0].(*int32) = r.Balance
	return nil
}
