package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	PassHash string `gorm:"not null"`
	IsAdmin  bool   `gorm:"not null;default:false"`

	Posts []Post `gorm:"foreignKey:UserID"`
	Favs  []UserFav
}
