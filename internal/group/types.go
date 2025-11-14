package group

import "github.com/jackc/pgx/v5/pgtype"

type GroupCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description" validate:"required"`
}

type GroupUpdateRequest struct {
	Name        *string `json:"name" validate:"min=3"`
	Description *string `json:"description"`
}

type GroupResponse struct {
	Id          string             `json:"id"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}
