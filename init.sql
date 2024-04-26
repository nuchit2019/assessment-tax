CREATE TABLE IF NOT EXISTS TaxBracket (
    BracketID SERIAL PRIMARY KEY,
    BracketName VARCHAR(255),
    MinIncome DECIMAL,
    MaxIncome DECIMAL,
    TaxRate DECIMAL
);

INSERT INTO TaxBracket (BracketName,MinIncome, MaxIncome, TaxRate) VALUES
('0.00-150,000.00',0, 150000, 0),
('150,001.00-500,000',150001, 500000, 10),
('500,001.00-1,000,000',500001, 1000000, 15),
('1,000,001.00-2,000,000',1000001, 2000000, 20),
('2,000,001.00 ขึ้นไป',2000001, NULL, 35);


CREATE TABLE IF NOT EXISTS Deduction (
    DeductionID SERIAL PRIMARY KEY,
    DeductionType VARCHAR(50) UNIQUE,
    DeductionAmount DECIMAL CHECK (DeductionAmount >= 0 AND DeductionAmount <= 100000)
);

INSERT INTO Deduction (DeductionType, DeductionAmount) VALUES
('personal', 60000),
('donation', 100000),
('k-receipt', 0);