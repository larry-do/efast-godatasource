package transaction

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"godatasource"
	"gorm.io/gorm"
)

func ExecuteInTransaction(f func(*gorm.DB) error) error {
	return ExecuteInTransactionWithDB(godatasource.DefaultConnection(), f)
}

func ExecuteInTransactionWithDB(db *gorm.DB, f func(*gorm.DB) error) error {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	txId, _ := uuid.NewUUID()
	log.Debug().Any("txId", txId).Msg("Start transaction")
	err := f(tx)
	if err != nil {
		tx.Rollback()
		log.Warn().Any("txId", txId).Msg("Got error. Rollback transaction")
	} else {
		tx.Commit()
		log.Debug().Any("txId", txId).Msg("Transaction committed")
	}
	return err
}
