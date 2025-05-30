CREATE TABLE payments(
    id VARCHAR(50) PRIMARY KEY,
    loan_id VARCHAR(50) REFERENCES loans(id),
    amount NUMERIC NOT NULL,
    start_at TIMESTAMP WITH TIME ZONE NOT NULL,
    end_at TIMESTAMP WITH TIME ZONE NOT NULL,
    paid_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    created_by INTEGER DEFAULT NULL,
    updated_by INTEGER DEFAULT NULL,
    deleted_by INTEGER DEFAULT NULL
)
CREATE INDEX IDX_loan_id_start_at ON payments(loan_id, start_at);