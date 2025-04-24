package migrations

import (
	"app/models"
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// 20250424140453-CreateRoute

func init() {
	goose.AddMigrationContext(up20250424140453, down20250424140453)
}
func createModel20250424140453 () interface{} {
	type Route struct{
		gorm.Model
		Name				string						`gorm:"type:varchar(100);"`
		Methods			string						`gorm:"type:varchar(100);"`
		Path				string						`gorm:"type:varchar(300);"`
		Header			datatypes.JSON		`gorm:"type:varchar(300);"`		// {[field]: {value: string, isReg: bool}}
		StripPath		bool
		UpStreamId	uint
		UpStream		models.UpStream
		Enable			bool
	}
	return &Route{}
}
func up20250424140453(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250424140453()
		createTable(migrator, model)
	})
}

func down20250424140453(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250424140453()
		dropTable(migrator, model)
	})
}
