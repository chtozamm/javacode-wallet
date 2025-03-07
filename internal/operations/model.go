package operations

const (
	Deposit  = "deposit"
	Withdraw = "withdraw"
)

type Operation struct {
	OperationType string `json:"operation_type"`
	Amount        int32  `json:"amount"`
}
