package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// 20250424135937-CreateServer

func init() {
	goose.AddMigrationContext(up20250424135937, down20250424135937)
}
func createModel20250424135937 () interface{} {
	type Server struct{
		gorm.Model
		Name				string						`gorm:"type:varchar(100);"`
		Host				string						`gorm:"type:varchar(100);"`
		Port				uint
		EnableSSL		bool
		Enable			bool
	}
	return &Server{}
}
func up20250424135937(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250424135937()
		createTable(migrator, model)
	})
}

func down20250424135937(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250424135937()
		dropTable(migrator, model)
	})
}
