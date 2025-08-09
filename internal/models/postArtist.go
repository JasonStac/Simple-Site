package models

type PostArtist struct {
	PostID   uint `gorm:"primaryKey"`
	ArtistID uint `gorm:"primaryKey"`

	Post   Post   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	Artist Artist `gorm:"foreignKey:ArtistID;constraint:OnDelete:CASCADE"`
}
