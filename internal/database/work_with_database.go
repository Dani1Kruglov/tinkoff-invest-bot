package database

import (
	"errors"
	"gorm.io/gorm"
	"tinkoff-investment-bot/internal/model"
)

func AddUser(db *gorm.DB, user *model.User) error {
	result := db.FirstOrCreate(user, model.User{TelegramID: user.TelegramID, Token: user.Token})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func AddAccount(db *gorm.DB, account *model.Account, userID uint) error {
	if userID != 0 {
		account.UserID = userID
		result := db.FirstOrCreate(account, model.Account{AccountID: account.AccountID, UserID: userID})
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		return errors.New("error select user by telegram id")
	}
}

func AddShare(db *gorm.DB, share *model.Share, userID uint, price float32) error {
	result := db.FirstOrCreate(share, model.Share{UID: share.UID})
	if result.Error != nil {
		return result.Error
	}

	result = db.FirstOrCreate(&model.UserFavoriteShare{
		UserID:    userID,
		ShareID:   share.ID,
		LastPrice: price,
	}, model.UserFavoriteShare{
		UserID:  userID,
		ShareID: share.ID,
	})
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

func GetAccountIDByTelegramID(db *gorm.DB, telegramID string) model.Account {
	var account model.Account
	db.Table("accounts").Select("*").Joins("left join users on accounts.user_id = users.id").Where("users.telegram_id = ?", telegramID).Scan(&account)
	return account
}
