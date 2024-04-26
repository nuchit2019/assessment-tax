package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

type Store interface {
	GetTaxBrackets() ([]model.TaxBracket, error)
	GetTaxLevel() ([]model.TaxLevel, error)
	UpdatePersonalDeduction(amount float64,deductType string) error
}

type Controller struct {
	store Store
}

func New(db Store) *Controller {
	return &Controller{store: db}
}

func (c *Controller) TaxCalculate(ctx echo.Context) error {
	var req model.TaxRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
	}

	taxableIncome := req.TotalIncome - 60000.0 // Apply standard deduction

	donationAllowance := 0.0
	kReceiptAllowance := 0.0
	for _, allowance := range req.Allowance {
		if allowance.Type == model.AllowanceTypeDonation {
			//taxableIncome -= allowance.Amount
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
	tax -= req.WHT

	if ctx.QueryParam("detail") == "true" {
		taxLevels := c.calculateTaxLevels(taxableIncome, tax)
		return ctx.JSON(http.StatusOK, model.TaxDetailResponse{Tax: tax, TaxLevel: taxLevels})
	}

	return ctx.JSON(http.StatusOK, model.TaxResponse{Tax: tax})
}

func (c *Controller) calculateTaxLevels(taxableIncome, tax float64) []model.TaxLevel {

	taxLevels, err := c.store.GetTaxLevel()
	if err != nil {
		//TODD Handle error appropriately (logging, internal error response)
		return nil
	}

	// Calculate tax levels based on taxable income
	if taxableIncome <= 150000 {
		taxLevels[0].Tax = 0.0
	} else if taxableIncome <= 500000 {
		taxLevels[1].Tax = tax //19000.0
	} else {
		// For income above 500,000, tax is 0 according to the new tax bracket structure
		taxLevels[1].Tax = tax //19000.0
	}

	return taxLevels
}

func (c *Controller) calculateTax(taxableIncome float64) float64 {
	var tax float64
	taxBrackets, err := c.store.GetTaxBrackets()
	if err != nil {
		// TODO Handle error appropriately (logging, internal error response)
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

		tax += taxableAmount * bracket.TaxRate
		taxableIncome -= taxableAmount
	}

	return tax
}
