-- name: GetWallets :many
SELECT * FROM wallets ORDER BY created_at;

-- name: GetBalance :one
SELECT balance FROM wallets 
WHERE id = $1 LIMIT 1;

-- name: CreateWallet :one
INSERT INTO wallets (id)
VALUES (
	gen_random_uuid()
)
RETURNING id;

-- name: AddOperation :exec
INSERT INTO operations (id, wallet_id, operation_type, amount)
VALUES (
	gen_random_uuid(),
	$1,
	$2,
	$3
);

-- name: UpdateWallet :exec
UPDATE wallets SET balance = $1, updated_at = NOW()
WHERE id = $2;

-- name: DeleteWallet :exec
DELETE FROM wallets WHERE id = $1;
