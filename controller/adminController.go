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

	if deductType == "personal" && req.Amount < 10000 {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "ค่าลดหย่อนส่วนตัวต้องมีค่ามากกว่า 10,000 บาท"})
	} 

	if deductType == "k-receipt" && req.Amount <= 0 {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "ค่าลด k-receipt ต้องมีค่ามากกว่า 0 บาท"})
	}
	
	if deductType == "k-receipt" && req.Amount > 100000 {
		return ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "แอดมิน สามารถกำหนด k-receipt สูงสุดได้ แต่ไม่เกิน 100,000 บาท"})
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
