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
	req := new(model.TaxRequest) 

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
	}

	taxableIncome, err := c.calculateTaxableIncome(req)
	if err != nil {
		return err
	}

	tax := c.calculateTax(taxableIncome) - req.WHT

	return ctx.JSON(http.StatusOK, model.TaxResponse{Tax: tax})
}

func (c *Controller) calculateTaxableIncome(req *model.TaxRequest) (float64, error) {
	taxableIncome := req.TotalIncome - 60000.0

	for _, allowance := range req.Allowance {
		switch allowance.Type {
		case model.AllowanceTypeDonation:
			taxableIncome -= allowance.Amount
			
		case model.AllowanceTypeKReceipt:
			// taxableIncome -= allowance.Amount
			// TODO: Implement K-Receipt logic
			return 0, echo.NewHTTPError(http.StatusBadRequest, "TODO: Implement K-Receipt logic")

		default:
			return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid allowance type")
		}
	}

	return taxableIncome, nil
}

func (c *Controller) calculateTax(taxableIncome float64) float64 {
	taxBrackets, err := c.store.GetTaxBrackets()
	if err != nil {
		// TODO: Handle error appropriately (logging, internal error response)
		return 0 
	}
	 

	var tax float64
	remainingIncome := taxableIncome

	for _, bracket := range taxBrackets {
		if remainingIncome <= 0 {
			break
		}

		taxableAmount := remainingIncome
		if taxableAmount <= bracket.MaxIncome {
			taxableAmount = remainingIncome
		} else {
			taxableAmount = bracket.MaxIncome
		}

		tax += taxableAmount * bracket.TaxRate
		remainingIncome -= taxableAmount
	}

	return tax
}