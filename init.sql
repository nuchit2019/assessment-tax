CREATE TABLE IF NOT EXISTS tax_bracket (
    id SERIAL PRIMARY KEY, bracket_name VARCHAR(255), min_income DECIMAL, max_income DECIMAL, tax_rate DECIMAL
);

INSERT INTO tax_bracket (id,bracket_name,min_income, max_income, tax_rate) VALUES
	 (1,'0-150,000',0,150000,0),
	 (2,'150,001-500,000',150000,500000,10),
	 (3,'500,001-1,000,000',500000,1000000,15),
	 (4,'1,000,001-2,000,000',1000000,2000000,20),
	 (5,'2,000,001 ขึ้นไป',2000000,'infinity'::numeric,35);

CREATE TABLE IF NOT EXISTS deduction (
    id SERIAL PRIMARY KEY, deduction_type VARCHAR(50) UNIQUE, deduction_amount DECIMAL CHECK (
        deduction_amount >= 0
        AND deduction_amount <= 100000
    )
);

INSERT INTO
    deduction (
        deduction_type, deduction_amount
    )
VALUES ('personal', 60000),
    ('donation', 100000),
    ('k-receipt', 50000);