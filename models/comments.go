package models

import "time"

type Comment struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	UserId    string    `json:"user_id"`
	PostId    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
