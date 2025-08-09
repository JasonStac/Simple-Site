package models

type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	MediaType MediaType `gorm:"type:media_types;not null"`
	Filename  string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"` // foreign key to user

	Owner   User `gorm:"constraint:OnDelete:CASCADE"`
	Favs    []UserFav
	Tags    []PostTag
	Artists []PostArtist
}
