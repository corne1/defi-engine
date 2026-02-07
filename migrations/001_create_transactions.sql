CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    hash TEXT,
    state TEXT NOT NULL,
    block_number BIGINT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
