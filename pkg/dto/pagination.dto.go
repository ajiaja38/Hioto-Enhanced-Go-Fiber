package dto

type PaginationRequest struct {
	Page   int    `json:"page" query:"Page" validate:"required,min=1"`
	Limit  int    `json:"limit" query:"Limit" validate:"required,min=1,max=100"`
	Search string `json:"search" query:"search" validate:"omitempty"`
}

type GetRulesPagination struct {
	PaginationRequest
	GuidSensor string `json:"guid_sensor" query:"guid_sensor" validate:"omitempty"`
}

type MetaPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
	TotalData int `json:"totalData"`
}
