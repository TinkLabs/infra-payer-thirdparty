package valoot

type TransactionResp struct {
	Data Transaction `json:"data"`
	Meta Meta        `json:"meta"`
}

func (r *TransactionResp) HasNext() bool {
	return r.Meta.Pagination.HasNext()
}

func (r *TransactionResp) GetNextPage() int {
	return r.Meta.Pagination.GetNextPage()
}

func (r *TransactionResp) HasPrevious() bool {
	return r.Meta.Pagination.HasPrevious()
}

func (r *TransactionResp) GetPreviousPage() int {
	return r.Meta.Pagination.GetPreviousPage()
}

type TransactionsResp struct {
	Data []Transaction `json:"data"`
	Meta Meta          `json:"meta"`
}

type PaymentParams struct {
	AppID     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type Transaction struct {
	Object        string         `json:"object"`
	ID            string         `json:"id"`
	Status        string         `json:"status"`
	Currency      string         `json:"currency"`
	Amount        float64        `json:"amount"`
	Message       string         `json:"message"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	ExpiredAt     *string        `json:"expired_at"`
	PaymentURL    *string        `json:"payment_url"`
	ImageURL      *string        `json:"image_url"`
	RealID        int            `json:"real_id"`
	PaymentParams *PaymentParams `json:"payment_params"`
}

type CreateTransactionParams struct {
	Currency        string  `json:"currency"`
	Amount          string  `json:"amount"`
	CallbackUrl     string  `json:"callback_url"`
	RedirectUrl     *string `json:"redirect_url,omitempty"`
	ProductName     string  `json:"product_name"`
	Service         string  `json:"service"`
	Wallet          string  `json:"wallet"`
	TransactionType *string `json:"type,omitempty"`
	AuthCode        *int    `json:"auth_code,omitempty"`
	Openid          *string `json:"openid,omitempty"`
	Extra           *string `json:"extra,omitempty"`
}

type ListTransactionParams struct {
	Id                  *string
	TransactionType     *string
	Status              *string
	ProviderReferenceId *string
	ProviderTradeId     *string
	ProviderOpenid      *string
	ProductName         *string
	CreatedAt           *string
	Limit               *int
	Page                *int
}
