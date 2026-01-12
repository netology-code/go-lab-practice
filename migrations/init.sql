CREATE TYPE transaction_type AS ENUM ('deposit', 'withdraw');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    balance NUMERIC(10,2) DEFAULT 0.00
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    type transaction_type NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    processed BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
