package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// 20250424140530-CreateUpstream

func init() {
	goose.AddMigrationContext(up20250424140530, down20250424140530)
}
func createModel20250424140530 () interface{} {
	type UpStream struct{
		gorm.Model
		Name				string						`gorm:"type:varchar(100);"`
		Scheme			string						`gorm:"type:varchar(20);"`		// eg: https://
		Address			string						`gorm:"type:varchar(300);"`		// host:port
	}
	return &UpStream{}
}
func up20250424140530(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250424140530()
		createTable(migrator, model)
	})
}

func down20250424140530(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250424140530()
		dropTable(migrator, model)
	})
}
