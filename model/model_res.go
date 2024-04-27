package model

type TaxResponse struct {
	Tax float64 `json:"tax"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type TaxDetailResponse struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
}

type TaxLevelRes struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxLevelModel struct {
	Id         int     `postgres:"id" json:"id"`
	Level      int     `postgres:"level" json:"level"`
	Label      string  `postgres:"label" json:"label"`
	MinAmount  float64 `postgres:"min_amount" json:"minAmount"`
	MaxAmount  float64 `postgres:"max_amount" json:"maxAmount"`
	TaxPercent int     `postgres:"tax_percent" json:"taxPercent"`
	Tax   float64 `json:"tax"`
}


type PersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type KreceiptDeductionResponse struct {
	KreceiptDeduction float64 `json:"kReceipt"`
}

type TaxCsvDetail struct {
	TotalIncome float64 `json:"totalIncome"`
    Tax         float64 `json:"tax"`
    TaxRefund   float64 `json:"taxRefund,omitempty"`
}

type TaxCsvDetailResponse struct {
	Taxes []TaxCsvDetail `json:"taxes"`
}
