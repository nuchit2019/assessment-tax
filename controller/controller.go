package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	// TODO: Implement Store interface
}

type controller struct {
	store Store
}

func New(db Store) *controller {
	return &controller{
		store: db,
	}
}

type Err struct {
	Message string `json:"message"`
}

type TaxRequest struct {
	TotalIncome float64 `json:"totalIncome" example:"500000.0"`
}

type TaxResponse struct {
	Tax float64 `json:"tax"`
}

type TaxBracket struct {
	MaxIncome float64 // Upper bound of the bracket
	TaxRate   float64 // Tax rate for this bracket
}

var taxBracket = []TaxBracket{
	{MaxIncome: 150000, TaxRate: 0},
	{MaxIncome: 500000, TaxRate: 0.10},
	{MaxIncome: 1000000, TaxRate: 0.15},
	{MaxIncome: 2000000, TaxRate: 0.20},
	{TaxRate: 0.35}, // Default/catch-all for incomes above 2000000

}

func (c *controller) TaxCalculate(ctx echo.Context) error {
	var req TaxRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Err{Message: "invalid request body"})
	}

	tax := calculateTax(req.TotalIncome)
	return ctx.JSON(http.StatusOK, TaxResponse{Tax: tax})
}

func calculateTax(income float64) float64 {
	tax := 0.0
	remainingIncome := income - 60000

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
