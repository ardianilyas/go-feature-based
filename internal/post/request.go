package post

import "github.com/google/uuid"

type CreatePostRequest struct {
	Title 		string 		`json:"title" binding:"required"`
	Content 	string 		`json:"content" binding:"required"`
	CategoryID 	uuid.UUID 	`json:"category_id" binding:"required"`
}

type UpdatePostRequest struct {
	Title 		string 		`json:"title" binding:"required"`
	Content 	string 		`json:"content" binding:"required"`
	CategoryID 	uuid.UUID 	`json:"category_id" binding:"required"`
}