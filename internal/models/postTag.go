package models

type PostTag struct {
	PostID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Tag  Tag  `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE"`
}
