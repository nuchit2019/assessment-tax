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

	err:=c.store.UpdatePersonalDeduction(req.Amount)
	if err!=nil{
		return ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "failed to update personal deduction"})	
	}
 
	return ctx.JSON(http.StatusOK, model.PersonalDeductionResponse{PersonalDeduction: req.Amount})
 
}

