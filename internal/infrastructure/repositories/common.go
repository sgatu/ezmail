package repositories

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func upsert(
	model interface{},
	ctx context.Context,
	db *bun.DB,
) error {
	result, err := db.NewUpdate().Model(model).WherePK().ExcludeColumn("id", "created").Exec(ctx)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 0 {
		fmt.Println("bun upsert: update")
		return nil
	}
	_, err = db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}
	fmt.Println("bun upsert: insert")
	return nil
}
