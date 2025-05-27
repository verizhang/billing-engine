CREATE TABLE loans(
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    amount NUMERIC NOT NULL,
    interest NUMERIC NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    created_by VARCHAR(50) DEFAULT NULL,
    updated_by VARCHAR(50) DEFAULT NULL,
    deleted_by VARCHAR(50) DEFAULT NULL
);
CREATE INDEX IDX_user_id_is_active ON loans(user_id, is_active);