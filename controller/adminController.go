package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

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
