package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

func (c *Controller) UpdatePersonalDeductionController(ctx echo.Context) error {

	var req model.PersonalDeductionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
	}

	deductType := ctx.Param("deductType")
	if deductType == "" {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "deductType Parameter is required"})
	}

	err := c.store.UpdatePersonalDeduction(req.Amount, deductType)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "failed to update personal deduction ERROR: " + err.Error()})
	}

	var responseBody interface{}
	switch deductType {
	case "personal":
		responseBody = model.PersonalDeductionResponse{PersonalDeduction: req.Amount}
	case "k-receipt":
		responseBody = model.KreceiptDeductionResponse{KreceiptDeduction: req.Amount}
	default:
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid deductType"})
	}

	return ctx.JSON(http.StatusOK, responseBody)
}
