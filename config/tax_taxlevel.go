package config

import (
	"fmt"
	"github.com/nuchit2019/assessment-tax/model"
)

func (c *Config) GetTaxLevel() ([]model.TaxLevel, error) {
	query := `
		SELECT bracket_name, 0 as tax_rate
		FROM tax_bracket
		ORDER BY id;
	`

	rows, err := c.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tax levels: %w", err)
	}
	defer rows.Close()

	var taxLevels []model.TaxLevel

	for rows.Next() {
		var level string
		var taxRate float64
		if err := rows.Scan(&level, &taxRate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		taxLevels = append(taxLevels, model.TaxLevel{Level: level, Tax: taxRate})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return taxLevels, nil
}