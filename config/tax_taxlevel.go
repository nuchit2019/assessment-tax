
package config

import "github.com/nuchit2019/assessment-tax/model"

func (c *Config) GetTaxLevel() ([]model.TaxLevel, error) {
	taxLevels := []model.TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 0.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}
	return taxLevels, nil
}