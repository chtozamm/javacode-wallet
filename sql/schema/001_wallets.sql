-- +goose Up
CREATE TABLE wallets(
	id UUID PRIMARY KEY,
	balance INTEGER NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE operations(
	id UUID PRIMARY KEY,
	wallet_id UUID NOT NULL,
	operation_type TEXT NOT NULL CHECK (operation_type IN ('deposit', 'withdraw')),
	amount INTEGER NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_wallet_id FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE wallets;
DROP TABLE operations;
