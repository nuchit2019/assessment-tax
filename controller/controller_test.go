package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nuchit2019/assessment-tax/model"
)

type MockStore struct{}

func (m *MockStore) GetTaxBrackets() ([]model.TaxBracket, error) {
	taxBrackets := []model.TaxBracket{
		{MinIncome: 0, MaxIncome: 150000, TaxRate: 0},
		{MinIncome: 150001, MaxIncome: 500000, TaxRate: 10},
		{MinIncome: 500001, MaxIncome: 1000000, TaxRate: 15},
		{MinIncome: 1000001, MaxIncome: 2000000, TaxRate: 20},
		{MinIncome: 2000001, MaxIncome: 0, TaxRate: 35},
	}
	return taxBrackets, nil
}

func (m *MockStore) GetTaxLevel() ([]model.TaxLevelModel, error) {
	taxLevels := []model.TaxLevelModel{
		{Level: 1, Label: "0-150,000", MinAmount: 0, MaxAmount: 150000, TaxPercent: 0},
		{Level: 2, Label: "150,001-500,000", MinAmount: 150001, MaxAmount: 500000, TaxPercent: 10},
		{Level: 3, Label: "500,001-1,000,000", MinAmount: 500001, MaxAmount: 1000000, TaxPercent: 15},
		{Level: 4, Label: "1,000,001-2,000,000", MinAmount: 1000001, MaxAmount: 2000000, TaxPercent: 20},
		{Level: 5, Label: "2,000,001 ขึ้นไป", MinAmount: 2000001, MaxAmount: 0, TaxPercent: 35},
	}
	return taxLevels, nil
}

func (m *MockStore) UpdatePersonalDeduction(amount float64, deductType string) error {
	fmt.Printf("Mock UpdatePersonalDeduction called with amount: %f, deductType: %s\n", amount, deductType)
	return nil
}

func (m *MockStore) GetDeduction() ([]model.Deduction, error) {
	deductions := []model.Deduction{
		{Id: 1, DeductionType: "personal", DeductionAmount: 5000},
		{Id: 2, DeductionType: "donation", DeductionAmount: 2000},
	}
	return deductions, nil
}

func TestTaxCalculate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Setup
		e := echo.New()
		reqBody := `{"totalIncome":100000,"allowance":[{"type":"donation","amount":50000},{"type":"kReceipt","amount":30000}]}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock store
		store := &MockStore{}
		ctrl := New(store)

		// Execute
		err := ctrl.TaxCalculate(c)

		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("I want to calculate my tax, totalIncome: 500000.0 return tax: 29000.0", func(t *testing.T) {
		// Setup
		e := echo.New()
		reqBody := `{"totalIncome":500000,"wht":0,"allowances":[{"allowanceType":"donation","amount":0}]}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock store
		store := &MockStore{}
		ctrl := New(store)

		// Execute
		err := ctrl.TaxCalculate(c)
 
		// Assert
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expected := `{"tax":29000}`
		result := rec.Body.String()

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, rec.Code)
		}

		// Parse expected JSON string
		var expectedTax model.TaxResponse
		if err := json.Unmarshal([]byte(expected), &expectedTax); err != nil {
			t.Errorf("failed to parse expected JSON: %v", err)
		}

		// Parse result JSON string
		var resultTax model.TaxResponse
		if err := json.Unmarshal([]byte(result), &resultTax); err != nil {
			t.Errorf("failed to parse result JSON: %v", err)
		}

		// Compare tax field
		if resultTax.Tax != expectedTax.Tax {
			t.Errorf("unexpected tax value: got %f, want %f", resultTax.Tax, expectedTax.Tax)
		}
 

	})

}
