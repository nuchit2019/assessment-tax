package controller

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

type Store interface {
	GetTaxBrackets() ([]model.TaxBracket, error)

	GetTaxLevel() ([]model.TaxLevelModel, error)
	UpdatePersonalDeduction(amount float64, deductType string) error
 
	GetDeduction()([]model.Deduction, error)
}

type Controller struct {
	store Store
}

func New(db Store) *Controller {
	return &Controller{store: db}
}

func (c *Controller) TaxCalculateController(ctx echo.Context) error {
	var req model.TaxRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
	}

	if len(req.Allowance) == 0 {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body allowances are empty"})
	}

	taxableIncome := req.TotalIncome - 60000.0 // Apply standard deduction
	donationAllowance := 0.0
	kReceiptAllowance := 0.0
	for _, allowance := range req.Allowance {
		if allowance.Type == model.AllowanceTypeDonation {
			if allowance.Amount > 100000 { // เงินบริจาคสามารถหย่อนได้สูงสุดเพียง 100,000 บาท
				donationAllowance += 100000
			} else {
				donationAllowance += allowance.Amount
			}

		} else if allowance.Type == model.AllowanceTypeKReceipt {
			if allowance.Amount > 50000 {
				kReceiptAllowance += 50000
			} else {
				kReceiptAllowance += allowance.Amount
			}
		}
	}

	taxableIncome -= (donationAllowance + kReceiptAllowance)

	tax := c.calculateTax(taxableIncome)

	return ctx.JSON(http.StatusOK, model.TaxResponse{Tax: tax})

}

func (c *Controller) calculateTaxLevels(taxableIncome, tax float64) []model.TaxLevelModel {

	taxLevels, err := c.store.GetTaxLevel()
	if err != nil {
		log.Printf("error getting tax levels: %v", err)
		return []model.TaxLevelModel{}
	}
	for i := range taxLevels {
		switch i {
		case 0:
			if taxableIncome <= 150000 {
				taxLevels[i].Tax = 0.00
			}
		case 1:
			if taxableIncome <= 500000 {
				if tax == 0 {
					taxLevels[i].Tax = 0.00
				} else {
					taxLevels[i].Tax = tax
				}
			}
		default:
			  taxLevels[i].Tax = tax
		}
	}

	return taxLevels
	  
}

func (c *Controller) calculateTax(taxableIncome float64) float64 {
	var tax float64
	taxBrackets, err := c.store.GetTaxBrackets()
	if err != nil {
		log.Printf("error getting tax brackets: %v", err)
		return 0
	}

	for _, bracket := range taxBrackets {
		if taxableIncome <= 0 {
			break
		}

		taxableAmount := taxableIncome
		if taxableAmount > bracket.MaxIncome {
			taxableAmount = bracket.MaxIncome
		}

		tax += taxableAmount * (bracket.TaxRate/100.0)
		taxableIncome -= taxableAmount
	}

	return tax
}


//===================AdminController=====================

type DeductionValidationError struct {
	Field   string
	Message string
}

func (e *DeductionValidationError) Error() string {
	return fmt.Sprintf("Field: %s, Message: %s", e.Field, e.Message)
}

func (c *Controller) UpdatePersonalDeductionController(ctx echo.Context) error {
	deductType := ctx.Param("deductType")
	if deductType == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "deductType Parameter is required"})
	}

	var req model.PersonalDeductionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
	}

	if err := validateDeduction(deductType, req.Amount); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
	}

	err := c.store.UpdatePersonalDeduction(req.Amount, deductType)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			model.ErrorResponse{Message: "failed to update personal deduction ERROR: " + err.Error()})
	}

	responseBody := constructResponse(deductType, req.Amount)

	return ctx.JSON(http.StatusOK, responseBody)

}

func validateDeduction(deductType string, amount float64) error {
	if deductType == "personal" && amount < 10000 {
		return &DeductionValidationError{Field: "amount", Message: "ค่าลดหย่อนส่วนตัวต้องมีค่ามากกว่า 10,000 บาท"}
	}
	if deductType == "k-receipt" && amount <= 0 {
		return &DeductionValidationError{Field: "amount", Message: "ค่าลด k-receipt ต้องมีค่ามากกว่า 0 บาท"}
	}
	if deductType == "k-receipt" && amount > 100000 {
		return &DeductionValidationError{Field: "amount", Message: "แอดมิน สามารถกำหนด k-receipt สูงสุดได้ แต่ไม่เกิน 100,000 บาท"}
	}

	return nil
}

func constructResponse(deductType string, amount float64) interface{} {
	switch deductType {
	case "personal":
		return model.PersonalDeductionResponse{PersonalDeduction: amount}
	case "k-receipt":
		return model.KreceiptDeductionResponse{KreceiptDeduction: amount}

	default:
		return nil
	}
}


//===================CSV_controller=====================

func (c *Controller) TaxCalculateFormCsv(ctx echo.Context) error {
	taxFile, err := ctx.FormFile("taxFile")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to retrieve tax file: " + err.Error()})
	}

	fileCsv, err := taxFile.Open()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to open tax file: " + err.Error()})
	}
	defer fileCsv.Close()

	data, err := io.ReadAll(fileCsv)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to read tax file: " + err.Error()})
	}

	deducts, err := c.store.GetDeduction()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve deductions: " + err.Error()})
	}
	deductType := DeductionTypeMap(deducts)	
	personalDeduction := personalType(deductType)
	donateDeduction := donationType(deductType)

	var taxCsv []model.TaxCsv
	err = gocsv.UnmarshalBytes(data, &taxCsv)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to parse tax CSV: " + err.Error()})
	}

	var taxes []model.TaxCsvDetail
	for _, t := range taxCsv {

		errs := validateTaxCsv(t)
		if len(errs) > 0 {
			return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid tax data: " + errs[0].Field+": "+errs[0].Message})
		}

		var tax float64
		deduct := 0.0

		totalIncome := t.TotalIncome
		wht := t.Wht
		donation := t.Donation

		if donation > donateDeduction {
			deduct = deduct + donateDeduction
		} else {
			deduct = deduct + donation
		}

		netIncome := (totalIncome - personalDeduction) - deduct
		allLevels, err := c.store.GetTaxLevel()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "failed to get tax level ERROR: " + err.Error()})
		}

		for _, t := range allLevels {
			eachtax := calcTaxByLevel(t, netIncome)
			tax += eachtax
		}

		tax = tax - wht

		if tax < 0 {
			taxes = append(taxes, model.TaxCsvDetail{
				TotalIncome: t.TotalIncome,
				Tax:         0.0,
				TaxRefund:   math.Abs(tax),
			})
		} else {
			taxes = append(taxes, model.TaxCsvDetail{
				TotalIncome: t.TotalIncome,
				Tax:         tax,
			})
		}
	}

	res := model.TaxCsvDetailResponse{
		Taxes: taxes,
	}

	return ctx.JSON(http.StatusOK, res)
}

func DeductionTypeMap(deducts []model.Deduction) map[string]float64 {
	data := make(map[string]float64)
	for _, item := range deducts {
		data[item.DeductionType] = item.DeductionAmount
	}
	return data
}

func personalType(item map[string]float64) float64 {
	return item["personal"]
}

func donationType(item map[string]float64) float64 {
	return item["donation"]
}
 
func validateTaxCsv(t model.TaxCsv) []model.ValidateErr {
	var errs []model.ValidateErr

	if t.TotalIncome < 0 {
		errs = append(errs, model.ValidateErr{
			Field:   "totalIncome",
			Message: "Total income must be more than 0",
		})
	}

	if t.Wht < 0 {
		errs = append(errs, model.ValidateErr{
			Field:   "wht",
			Message: "Withholding tax must be more than 0",
		})
	}

	if t.Wht > t.TotalIncome {
		errs = append(errs, model.ValidateErr{
			Field:   "wht",
			Message: "Withholding tax must be less than total income",
		})
	}

	if t.Donation < 0 {
		errs = append(errs, model.ValidateErr{
			Field:   "donation",
			Message: "Donation must be more than 0",
		})
	}

	return errs
}

func calcTaxByLevel(level model.TaxLevelModel, income float64) float64 {
	if income <= level.MinAmount {
		return 0
	}
	if income <= level.MaxAmount {
		return (income-level.MinAmount) * float64(level.TaxPercent) / 100
	}
	return (level.MaxAmount-level.MinAmount) * float64(level.TaxPercent) / 100
}
