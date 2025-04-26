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
// 20250426091214-AddSortToRoute

func init() {
	goose.AddMigrationContext(up20250426091214, down20250426091214)
}
func createModel20250426091214 () interface{} {
	type Route struct{
		gorm.Model
		Name				string						`gorm:"type:varchar(100);"`
		Methods			string						`gorm:"type:varchar(100);"`		// GET,POST,....
		Path				string						`gorm:"type:varchar(300);"`
		HeaderUp		datatypes.JSON			// []{key: string, value: string}
		HeaderDown	datatypes.JSON			// []{key: string, value: string}
		StripPath		bool
		UpStreamId	uint
		UpStream		models.UpStream
		Enable			bool
		ServerId		uint
		AuthenticationId	uint
		Authentication models.Authentication
		Sort	uint
	}
	return &Route{}
}
func up20250426091214(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250426091214()
		addTableColumn(migrator, model, "Sort")
	})
}

func down20250426091214(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250426091214()
		dropTableColumn(migrator, model, "Sort")
	})
}
