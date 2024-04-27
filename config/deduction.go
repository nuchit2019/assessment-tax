package config

import (
	// "database/sql"
	"fmt"

	"github.com/nuchit2019/assessment-tax/model"
)

func (c *Config) GetDeduction() ([]model.Deduction, error) {
	query := `
		SELECT deduction_type, deduction_amount
		FROM deduction
		ORDER BY id;
	`

	rows, err := c.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deduction: %w", err)
	}
	defer rows.Close()

	var deductions []model.Deduction
	for rows.Next() {
		var data model.Deduction
		err := rows.Scan(
			&data.DeductionType,
			&data.DeductionAmount,
		)
		if err != nil {
			return nil, err
		}
		deductions = append(deductions, model.Deduction{
			DeductionType:   data.DeductionType,
			DeductionAmount: data.DeductionAmount,
		})

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return deductions, nil
}
