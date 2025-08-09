package models

type UserFav struct {
	UserID uint `gorm:"primaryKey"`
	PostID uint `gorm:"primaryKey"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
