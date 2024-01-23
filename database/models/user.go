package models

type User struct {
	ID           int64
	UserName     string
	HashPassword string
	AvatarPath   string
}
