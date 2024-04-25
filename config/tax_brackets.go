
package config

import "github.com/nuchit2019/assessment-tax/model"

func (c *Config) GetTaxBrackets() ([]model.TaxBracket, error) {
	taxBrackets := []model.TaxBracket{
		{MaxIncome: 150000, TaxRate: 0},
		{MaxIncome: 500000, TaxRate: 0.10},
		{MaxIncome: 1000000, TaxRate: 0.15},
		{MaxIncome: 2000000, TaxRate: 0.20},
		{TaxRate: 0.35},
	}
	return taxBrackets, nil
}