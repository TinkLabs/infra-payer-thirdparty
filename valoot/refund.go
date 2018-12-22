package valoot

type RefundResp struct {
	Data Refund `json:"data"`
	Meta Meta   `json:"meta"`
}

type Refund struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Status    string `json:"status"`
	Amount    string `json:"amount"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
