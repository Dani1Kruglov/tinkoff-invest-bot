package database

import (
	"errors"
	"gorm.io/gorm"
	"tinkoff-investment-bot/internal/model"
)

func AddUser(db *gorm.DB, user *model.User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByTelegramID(db *gorm.DB, telegramID string) model.User {
	var user model.User
	db.Table("users").Select("*").Where("telegram_id = ?", telegramID).Scan(&user)
	return user
}

func AddAccount(db *gorm.DB, account *model.Account, userID uint) error {
	if userID != 0 {
		account.UserID = userID
		result := db.Create(account)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		return errors.New("error select user by telegram id")
	}
}

func GetAccountIDByTelegramID(db *gorm.DB, telegramID string) model.Account {
	var account model.Account
	db.Table("accounts").Select("*").Joins("left join users on accounts.user_id = users.id").Where("users.telegram_id = ?", telegramID).Scan(&account)
	return account
}
