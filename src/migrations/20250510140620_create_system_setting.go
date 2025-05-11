package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// 20250510140620-CreateSystemSetting

func init() {
	goose.AddMigrationContext(up20250510140620, down20250510140620)
}
func createModel20250510140620 () interface{} {
	type SystemSetting struct {
		gorm.Model
		Name				string						`gorm:"type:varchar(100);"`
		Value				datatypes.JSON		`gorm:"type:json;"`
	}
	return &SystemSetting{}
}
func up20250510140620(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250510140620()
		createTable(migrator, model)
	})
}

func down20250510140620(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250510140620()
		dropTable(migrator, model)
	})
}
