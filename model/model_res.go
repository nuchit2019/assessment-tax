package model


type TaxResponse struct {
	Tax float64 `json:"tax"`
}

type ErrorResponse  struct {
	Message string `json:"message"`
}