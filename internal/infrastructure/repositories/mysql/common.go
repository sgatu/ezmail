package mysql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/uptrace/bun"
)

func upsert(
	model interface{},
	ctx context.Context,
	db *bun.DB,
) error {
	slog.Debug(fmt.Sprintf("Upserting %+v", model))
	result, err := db.NewUpdate().Model(model).WherePK().ExcludeColumn("id", "created").Exec(ctx)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 0 {
		return nil
	}
	// necessary due to affectedRows = 0 if no field was changed
	res, err := db.NewSelect().Model(model).WherePK().Count(ctx)
	if err == nil && res >= 1 {
		return nil
	}
	_, err = db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
