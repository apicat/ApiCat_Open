package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184413",
		Migrate: func(tx *gorm.DB) error {

			type Tag struct {
				ID           uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
				ProjectID    string `gorm:"type:varchar(24);index;not null;comment:project id"`
				Name         string `gorm:"type:varchar(255);not null;comment:tag name"`
				DisplayOrder int    `gorm:"type:int(11);not null;default:0;comment:display order"`
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}

			if tx.Migrator().HasTable(&Tag{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&Tag{})
		},
	}

	MigrationHelper.Register(m)
}
