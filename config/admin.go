
package config

func (c *Config) UpdatePersonalDeduction(amount float64) error {
	_, err := c.Db.Exec("UPDATE deduction SET deduction_amount = $1 WHERE deduction_type = 'personal'", amount)
	return err
}