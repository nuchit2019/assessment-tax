package config

import (
	"fmt"
	"github.com/nuchit2019/assessment-tax/model"
)

// func (c *Config) GetTaxLevel() ([]model.TaxLevel, error) {
// 	taxLevels := []model.TaxLevel{
// 		{Level: "0-150,000", Tax: 0.0},
// 		{Level: "150,001-500,000", Tax: 0.0},
// 		{Level: "500,001-1,000,000", Tax: 0.0},
// 		{Level: "1,000,001-2,000,000", Tax: 0.0},
// 		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
// 	}
// 	return taxLevels, nil
// }

func (c *Config) GetTaxLevel() ([]model.TaxLevel, error) {
	// Query to select tax levels from the tax_bracket table
	query := `
		SELECT bracket_name, tax_rate
		FROM tax_bracket
		ORDER BY id;
	`

	// Execute the query to fetch tax levels
	rows, err := c.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tax levels: %w", err)
	}
	defer rows.Close()

	// Initialize a slice to store tax levels
	var taxLevels []model.TaxLevel

	// Iterate over the rows and populate the taxLevels slice
	for rows.Next() {
		var level string
		var taxRate float64
		if err := rows.Scan(&level, &taxRate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		taxLevels = append(taxLevels, model.TaxLevel{Level: level, Tax: taxRate})
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return taxLevels, nil
}