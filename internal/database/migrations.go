package database

import (
	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20240520_initial_payment_table",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Payment{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("payments")
			},
		},
		// Aqui você pode adicionar outras migrações no futuro
	})

	return m.Migrate()
}
