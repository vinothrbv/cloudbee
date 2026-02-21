package entity

import "time"

type Post struct {
	ID              int64 `gorm:"primaryKey"`
	Title           string
	Content         string
	Author          string
	Tags            []string
	PublicationDate time.Time
	CreatedAt       time.Time `gorm:"->"`
	UpdatedAt       time.Time `gorm:"->"`
}
