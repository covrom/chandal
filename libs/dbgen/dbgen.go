package dbgen

import (
	"context"
	"fmt"

	"github.com/covrom/chandal/libs/gjsonb"
	"github.com/jmoiron/sqlx"
)

func SelectAll[T any](ctx context.Context, tx *sqlx.Tx, table string,
	f func(context.Context, *T) error,
) error {
	rows, err := tx.QueryxContext(ctx,
		fmt.Sprintf(`SELECT to_jsonb(%[1]s) v FROM %[1]s`, table))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		v := gjsonb.JsonB[T]{}
		if err := rows.Scan(&v); err != nil {
			return err
		}
		if v.Val == nil {
			continue
		}
		if err := f(ctx, v.Val); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return rows.Err()
}

func Insert[T any](ctx context.Context, tx *sqlx.Tx, table string, values ...T) error {
	if _, err := tx.ExecContext(ctx,
		fmt.Sprintf(`INSERT INTO %[1]s SELECT *
			FROM json_populate_recordset(NULL::%[1]s, $1)`, table),
		gjsonb.JsonB[[]T]{Val: &values}); err != nil {
		return err
	}
	return nil
}
