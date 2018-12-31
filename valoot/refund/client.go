package refund

import (
	"strconv"

	"github.com/TinkLabs/go-services/valoot"
)

type Client struct {
	B valoot.Backend
}

func Refund(accessToken, transactionId string, amount float64) (*valoot.RefundResp, error) {
	return getC().Refund(accessToken, transactionId, amount)
}

func (c Client) Refund(accessToken, transactionId string, amount float64) (resp *valoot.RefundResp, err error) {
	content := map[string]string{}

	content["transaction_id"] = transactionId
	content["amount"] = strconv.FormatFloat(amount, 'f', -1, 64)

	path := "/v1/refunds"
	err = c.B.Call("POST", path, accessToken, nil, &content, &resp)
	return resp, err
}

func getC() Client {
	return Client{valoot.GetBackend(valoot.PublicBackend)}
}
