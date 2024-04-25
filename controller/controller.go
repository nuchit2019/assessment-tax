package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
) 
 
type Store interface {
	// TODO: Implement Store interface
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

	taxableIncome := req.TotalIncome - 60000.0

	for _, allowance := range req.Allowance {
		if allowance.Type == model.AllowanceTypeDonation {
			taxableIncome -= allowance.Amount
		}
	}

	tax := calculateTax(taxableIncome)

	tax -= req.WHT
	return ctx.JSON(http.StatusOK, model.TaxResponse{Tax: tax})
}

func calculateTax(income float64) float64 {
	taxBracket := []model.TaxBracket{
		{MaxIncome: 150000, TaxRate: 0},
		{MaxIncome: 500000, TaxRate: 0.10},
		{MaxIncome: 1000000, TaxRate: 0.15},
		{MaxIncome: 2000000, TaxRate: 0.20},
		{TaxRate: 0.35},
	}

	tax := 0.0
	remainingIncome := income

	for _, bracket := range taxBracket {
		if remainingIncome <= 0 {
			break
		}

		taxableAmount := remainingIncome
		if taxableAmount > bracket.MaxIncome {
			taxableAmount = bracket.MaxIncome
		}

		tax += taxableAmount * bracket.TaxRate
		remainingIncome -= taxableAmount
	}

	return tax
}