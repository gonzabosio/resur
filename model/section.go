package model

type Section struct {
	Id        int64  `json:"id,omitempty"`
	Title     string `json:"title" validate:"required"`
	ProjectID int64  `json:"project_id" validate:"required"`
}