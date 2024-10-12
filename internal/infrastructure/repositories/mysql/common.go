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
	slog.Info(fmt.Sprintf("Upserting %+v", model))
	result, err := db.NewUpdate().Model(model).WherePK().ExcludeColumn("id", "created").Exec(ctx)
	if err != nil {
		slog.Info(fmt.Sprintf("Upserting err %+v", err))
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		slog.Info(fmt.Sprintf("Upserting err aff %+v", err))
		return err
	}
	if affectedRows != 0 {
		return nil
	}
	_, err = db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
