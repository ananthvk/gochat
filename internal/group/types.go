package group

import (
	"time"
)

type GroupCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"required"`
}

type GroupUpdateRequest struct {
	Name        *string `json:"name" validate:"min=3"`
	Description *string `json:"description"`
}

type GroupResponse struct {
	Id          string    `json:"id"`
	OwnerId     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type MemberResponse struct {
	UsrId    string    `json:"usr_id"`
	JoinedAt time.Time `json:"joined_at"`
	Role     string    `json:"role"`
}

type MemberResponseWithName struct {
	UsrId    string    `json:"usr_id"`
	JoinedAt time.Time `json:"joined_at"`
	Role     string    `json:"role"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
}

type GroupListResponse struct {
	Id          string                    `json:"id"`
	OwnerId     string                    `json:"owner_id"`
	CreatedAt   time.Time                 `json:"created_at"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	LastMessage *GroupListMessageResponse `json:"last_message"`
}

type GroupListMessageResponse struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Type       string    `json:"type"`
	GrpId      string    `json:"group_id"`
	Content    string    `json:"content"`
	SenderId   string    `json:"sender_id"`
	SenderName string    `json:"sender_name"`
}
