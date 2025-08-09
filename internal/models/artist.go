package models

type Artist struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`

	Posts []PostArtist
}
