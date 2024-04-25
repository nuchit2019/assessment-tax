package model

import  "fmt"

type AllowanceType string

const (
	AllowanceTypeDonation AllowanceType = "donation"
	AllowanceTypeKReceipt AllowanceType = "k-receipt"
)

type Allowance struct {
	Type        AllowanceType 	`json:"allowanceType"`
	Amount      float64	  		`json:"amount"`
}

type TaxRequest struct {
	TotalIncome float64 `json:"totalIncome" example:"500000.0"`
	WHT	 		float64 `json:"wht"`
	Allowance   []Allowance `json:"allowances"`
}

type ErrInvalidAllowanceType  struct {
	Type AllowanceType
}

func (e *ErrInvalidAllowanceType) Error() string {
	return fmt.Sprintf("invalid allowance type: %s", e.Type)
}

