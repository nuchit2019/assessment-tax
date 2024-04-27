package config

import (
	"database/sql"
	"fmt"

	"github.com/nuchit2019/assessment-tax/model"
)

func (c *Config) GetTaxBrackets() ([]model.TaxBracket, error) {
	query := `
		SELECT min_income, max_income, tax_rate
		FROM tax_bracket
		ORDER BY id;
	`

	rows, err := c.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tax brackets: %w", err)
	}
	defer rows.Close()

	var taxBrackets []model.TaxBracket

	for rows.Next() {
		var minIncome, maxIncome sql.NullFloat64
		var taxRate sql.NullFloat64

		if err := rows.Scan(&minIncome, &maxIncome, &taxRate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		bracket := model.TaxBracket{}
		if minIncome.Valid {
			bracket.MinIncome = minIncome.Float64
		}
		if maxIncome.Valid {
			bracket.MaxIncome = maxIncome.Float64
		}
		if taxRate.Valid {
			bracket.TaxRate = taxRate.Float64
		}

		taxBrackets = append(taxBrackets, bracket)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return taxBrackets, nil
}
