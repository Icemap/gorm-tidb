package tidb

import (
	"errors"
	"gorm.io/gorm"
	"testing"
)

type Player struct {
	ID    uint `gorm:"primarykey;column:id"`
	Coins int  `gorm:"column:coins"`
	Goods int  `gorm:"column:goods"`
}

func TestNestedTxn(t *testing.T) {
	db, err := gorm.Open(Open("root:@tcp(127.0.0.1:4000)/test"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	// auto schema
	db.AutoMigrate(&Player{})

	// remove savepoint
	db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&Player{ID: 1, Coins: 1, Goods: 1})

		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&Player{ID: 2, Coins: 1, Goods: 1})
			return errors.New("rollback player2") // Rollback player2
		})

		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&Player{ID: 3, Coins: 1, Goods: 1})
			return nil
		})

		return nil
	})
}
