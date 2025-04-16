package dtos

type PaginationDTO struct {
	Page int `json:"page" default:"1"`
	Size int `json:"size" default:"10"`
}
