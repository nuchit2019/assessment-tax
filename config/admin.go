
package config

func (c *Config) UpdatePersonalDeduction(amount float64,deductType string) error {
	_, err := c.Db.Exec("UPDATE deduction SET deduction_amount = $1 WHERE deduction_type = $2", amount,deductType)
	return err
}