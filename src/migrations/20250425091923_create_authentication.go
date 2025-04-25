package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// 20250425091923-CreateAuthentication

func init() {
	goose.AddMigrationContext(up20250425091923, down20250425091923)
}
func createModel20250425091923 () interface{} {
	// TODO 修改对应的 ModelName
	type Authentication struct{
		gorm.Model
		Name				string						`gorm:"type:varchar(100);"`
		Username		string						`gorm:"type:varchar(100);"`
		Password		string						`gorm:"type:varchar(100);"`
		HashedPw		string						`gorm:"type:varchar(100);"`
	}
	return &Authentication{}
}
func up20250425091923(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250425091923()
		createTable(migrator, model)
	})
}

func down20250425091923(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250425091923()
		dropTable(migrator, model)
	})
}
