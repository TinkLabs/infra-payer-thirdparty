package transaction

import (
	"fmt"
	"net/url"

	"github.com/TinkLabs/payer-thirdparty/valoot"
)

type Client struct {
	B valoot.Backend
}

func CreateTransaction(accessToken string, params *valoot.CreateTransactionParams) (*valoot.TransactionResp, error) {
	return getC().CreateTransaction(accessToken, params)
}

func GetTransaction(accessToken string, id string) (*valoot.TransactionResp, error) {
	return getC().GetTransaction(accessToken, id)
}

func GetTransactions(accessToken string, params *valoot.ListTransactionParams) (*valoot.TransactionsResp, error) {
	return getC().GetTransactions(accessToken, params)
}

func (c Client) CreateTransaction(accessToken string, params *valoot.CreateTransactionParams) (resp *valoot.TransactionResp, err error) {
	path := "/v1/transactions"
	err = c.B.Call("POST", path, accessToken, nil, &params, &resp)
	return resp, err
}

func (c Client) GetTransaction(accessToken, transactionId string) (resp *valoot.TransactionResp, err error) {
	path := fmt.Sprintf("/v1/transactions/%s", transactionId)
	err = c.B.Call("GET", path, accessToken, nil, nil, &resp)
	return resp, err
}

func (c Client) GetTransactions(accessToken string, params *valoot.ListTransactionParams) (resp *valoot.TransactionsResp, err error) {
	v := url.Values{}

	if params.Id != nil {
		v.Add("id", fmt.Sprintf("%s", *params.Id))
	}
	if params.TransactionType != nil {
		v.Add("type", fmt.Sprintf("%s", *params.TransactionType))
	}
	if params.Status != nil {
		v.Add("status", fmt.Sprintf("%s", *params.Status))
	}
	if params.ProviderReferenceId != nil {
		v.Add("provider_reference_id", fmt.Sprintf("%s", *params.ProviderReferenceId))
	}
	if params.ProviderTradeId != nil {
		v.Add("provider_trade_id", fmt.Sprintf("%s", *params.ProviderTradeId))
	}
	if params.ProviderOpenid != nil {
		v.Add("provider_openid", fmt.Sprintf("%s", *params.ProviderOpenid))
	}
	if params.ProductName != nil {
		v.Add("product_name", fmt.Sprintf("%s", *params.ProductName))
	}
	if params.CreatedAt != nil {
		v.Add("created_at", fmt.Sprintf("%s", *params.CreatedAt))
	}
	if params.Limit != nil {
		v.Add("limit", fmt.Sprintf("%s", string(*params.Limit)))
	}
	if params.Page != nil {
		v.Add("page", fmt.Sprintf("%s", string(*params.Page)))
	}

	path := "/v1/transactions"
	err = c.B.Call("GET", path, accessToken, &v, nil, &resp)
	return resp, err
}

func getC() Client {
	return Client{valoot.GetBackend(valoot.APIBackend)}
}
