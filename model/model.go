package model

type TaxBracket struct {
	MaxIncome float64
	MinIncome float64
	TaxRate   float64
}

type Deduction struct {
	Id              int     `postgres:"id" json:"id"`
	DeductionType   string  `postgres:"deduction_type" json:"deduction_type"`
	DeductionAmount float64 `postgres:"deduction_amount" json:"deduction_amount"`
}

// TaxCsv represents the structure of CSV data for tax calculation.
type TaxCsv struct {
	TotalIncome float64 `csv:"totalIncome"`
	Wht         float64 `csv:"wht"`
	Donation    float64 `csv:"donation"`
}

type ValidateErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}