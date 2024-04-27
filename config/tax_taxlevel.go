package config

import (
	"fmt"
	"github.com/nuchit2019/assessment-tax/model"
)

func (c *Config) GetTaxLevel() ([]model.TaxLevelModel, error) {

	query := `
	SELECT id, bracket_name, min_income, max_income, tax_rate
	FROM tax_bracket
	ORDER BY id;
`

	rows, err := c.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tax levels: %w", err)
	}
	defer rows.Close()

	var levels []model.TaxLevelModel

	for rows.Next() {
		var l model.TaxLevelModel
		err := rows.Scan(
			&l.Level,
			&l.Label,
			&l.MinAmount,
			&l.MaxAmount,
			&l.TaxPercent,
		)

		if err != nil {
			return nil, err
		}
		levels = append(levels, model.TaxLevelModel{
			Level:      l.Level,
			Label:      l.Label,
			MinAmount:  l.MinAmount,
			MaxAmount:  l.MaxAmount,
			TaxPercent: l.TaxPercent,
		})
	}
	return levels, nil
}