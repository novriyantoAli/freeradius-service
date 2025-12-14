package database

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager struct {
	db *gorm.DB
}

func (tm *TransactionManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = WithTx(ctx, tx)
		return fn(ctx)
	})
}
