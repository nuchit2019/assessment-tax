package model


type TaxResponse struct {
	Tax float64 `json:"tax"`
}

type ErrorResponse  struct {
	Message string `json:"message"`
}

type TaxDetailResponse struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
}

// TaxLevel represents a tax level detail.
type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type PersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type KreceiptDeductionResponse struct {
	KreceiptDeduction float64 `json:"kReceipt"`
}
