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
// 20250425144410-AddFieldToRoute

func init() {
	goose.AddMigrationContext(up20250425144410, down20250425144410)
}
func createModel20250425144410 () interface{} {
	type Route struct{
		gorm.Model
		Name				string						`gorm:"type:varchar(100);"`
		Methods			string						`gorm:"type:varchar(100);"`		// GET,POST,....
		Path				string						`gorm:"type:varchar(300);"`
		Header			datatypes.JSON			// {[field]: {value: string, isReg: bool}}
		StripPath		bool
		UpStreamId	uint
		UpStream		models.UpStream
		Enable			bool
		ServerId		uint
		AuthenticationId	uint
		Authentication models.Authentication
	}
	return &Route{}
}
func up20250425144410(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250425144410()
		addTableColumn(migrator, model, "ServerId")
		addTableColumn(migrator, model, "AuthenticationId")
	})
}

func down20250425144410(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250425144410()
		dropTableColumn(migrator, model, "ServerId")
		dropTableColumn(migrator, model, "AuthenticationId")
	})
}
