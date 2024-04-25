package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	// TODO: Implement Store interface
}

type Controller  struct {
	store Store
}

func New(db Store) *Controller {
    return &Controller{store: db}
}

// type ErrorResponse  struct {
// 	Message string `json:"message"`
// }

// type TaxRequest struct {
// 	TotalIncome float64 `json:"totalIncome" example:"500000.0"`
// }

// type TaxResponse struct {
// 	Tax float64 `json:"tax"`
// }

// type TaxBracket struct {
// 	MaxIncome float64 
// 	TaxRate   float64
// }

var taxBracket = []model.TaxBracket{
	{MaxIncome: 150000, TaxRate: 0},
	{MaxIncome: 500000, TaxRate: 0.10},
	{MaxIncome: 1000000, TaxRate: 0.15},
	{MaxIncome: 2000000, TaxRate: 0.20},
	{TaxRate: 0.35},

}

func (c *Controller) TaxCalculate(ctx echo.Context) error {
	var req model.TaxRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
	}

	tax := calculateTax(req.TotalIncome)
	return ctx.JSON(http.StatusOK, model.TaxResponse{Tax: tax})
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
