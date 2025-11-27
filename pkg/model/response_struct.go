package model

type ResponseEntity[T any] struct {
	Code    int             `json:"code"`
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    T               `json:"data"`
	Meta    *MetaPagination `json:"meta,omitempty"`
}

type MetaPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
	TotalData int `json:"totalData"`
}

type ResponseError[T any] struct {
	ResponseEntity[T]
	Path string `json:"path"`
}
