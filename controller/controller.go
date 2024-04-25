package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

type Store interface {
	GetTaxBrackets() ([]model.TaxBracket, error)
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
	for _, allowance := range req.Allowance {
		if allowance.Type == model.AllowanceTypeDonation {
			//taxableIncome -= allowance.Amount
			if allowance.Amount > 100000 { // เงินบริจาคสามารถหย่อนได้สูงสุดเพียง 100,000 บาท
				donationAllowance += 100000
			} else {
				donationAllowance += allowance.Amount
			}

		}
	}

	taxableIncome -= donationAllowance

	tax := c.calculateTax(taxableIncome)
	tax -= req.WHT
	return ctx.JSON(http.StatusOK, model.TaxResponse{Tax: tax})
}

func (c *Controller) calculateTax(taxableIncome float64) float64 {
	var tax float64
	taxBrackets, err := c.store.GetTaxBrackets()
	if err != nil {
		// Handle error appropriately (logging, internal error response)
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
