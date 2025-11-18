package message

import "time"

// The content of a message can have atmost 4096 characters

type MessageCreateRequest struct {
	Type    string `json:"type" validate:"required,oneof=text"`
	Content string `json:"content" validate:"required,max=4096"`
}

type Cursor struct {
	Before    string `json:"before"`
	HasBefore bool   `json:"has_before"`
}

type MessageResponse struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
	GrpId     string    `json:"group_id"`
	Content   string    `json:"content"`
}

type MessagePaginationResponse struct {
	Cursor    Cursor          `json:"cursor"`
	Messsages MessageResponse `json:"messages"`
}
