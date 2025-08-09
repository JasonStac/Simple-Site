package models

type Session struct {
	Username  string `gorm:"primaryKey"`
	SessionID string `gorm:"uniqueIndex;not null"`
}
