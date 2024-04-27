package controller

import (
	"io"
	"math"
	"net/http" 

	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

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