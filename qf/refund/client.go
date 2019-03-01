package refund

import (
	"github.com/TinkLabs/payer-thirdparty/qf"
)

type Client struct {
	B       qf.Backend
	AppCode string
}

func New(sign string, params *qf.RefundParams) (*qf.Refund, error) {
	return getC().New(sign, params)
}

func (c Client) New(sign string, params *qf.RefundParams) (*qf.Refund, error) {
	path := "/v1/refund"
	refund := &qf.Refund{}

	err := c.B.Call("POST", path, c.AppCode, sign, nil, params, refund)
	return refund, err
}

func getC() Client {
	return Client{qf.GetBackend(qf.APIBackend), qf.AppCode}
}
