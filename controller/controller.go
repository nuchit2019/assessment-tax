package controller

import (
 
	"log"
	"net/http"

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

func (c *Controller) TaxCalculate(ctx echo.Context) error {
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
			// TODO Handle other cases if necessary
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
