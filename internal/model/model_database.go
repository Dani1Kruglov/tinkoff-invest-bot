package model

type User struct {
	ID         uint   `gorm:"primaryKey"`
	TelegramID string `gorm:"column:telegram_id;type:varchar(100)"`
	Token      string `gorm:"column:token;type:varchar(100)"`
}

type Account struct {
	AccountID string `gorm:"primaryKey;type:varchar(100)"`
	Name      string `gorm:"column:name;type:varchar(100)"`
	UserID    uint   `gorm:"column:user_id;index" json:"-"`
	User      User   `gorm:"foreignKey:UserID"`
}

type Share struct {
	UID       string `gorm:"primaryKey;type:varchar(200)"`
	Ticker    string `gorm:"column:ticker;type:varchar(10)"`
	Name      string `gorm:"column:name;type:varchar(100)"`
	FIGI      string `gorm:"column:figi;type:varchar(200)"`
	ClassCode string `gorm:"column:class_code;type:varchar(10)"`
}

type UserFavoriteShare struct {
	UserID    uint    `gorm:"column:user_id;index;type:int"`
	ShareID   string  `gorm:"column:share_id;index;type:varchar(200)"`
	LastPrice float32 `gorm:"column:last_price;type:float"`
	Share     Share   `gorm:"foreignKey:ShareID"`
	User      User    `gorm:"foreignKey:UserID"`
}
