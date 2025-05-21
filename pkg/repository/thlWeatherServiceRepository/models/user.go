package models

import (
	"time"
)

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey" swaggertype:"integer"`
	Email     string     `json:"email" gorm:"column:email" swaggertype:"string"`
	Password  string     `json:"password" gorm:"column:password" swaggertype:"string"`
	FullName  string     `json:"full_name" gorm:"column:full_name" swaggertype:"string"`
	Profile   string     `json:"profile" gorm:"column:profile" swaggertype:"string"`
	Likes     int64      `json:"likes" gorm:"column:likes;default:0" swaggertype:"integer"`
	Bookmark  int64      `json:"bookmark" gorm:"column:bookmark;default:0" swaggertype:"integer"`
	Fowlow    int64      `json:"follow" gorm:"column:follow;default:0" swaggertype:"integer"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at" swaggertype:"string" format:"date-time"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at" swaggertype:"string" format:"date-time"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at" swaggertype:"string" format:"date-time"`
}
