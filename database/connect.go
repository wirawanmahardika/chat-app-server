package database

import (
	"chatapp/database/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("storage/system.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Conversation{},
		&model.ConversationParticipant{},
		&model.Friendship{},
		&model.Message{},
	); err != nil {
		panic(err)
	}

	return db
}
