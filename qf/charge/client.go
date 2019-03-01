package charge

import (
	"github.com/TinkLabs/payer-thirdparty/qf"
)

type Client struct {
	B       qf.Backend
	AppCode string
}

// New creates a new charge.
func New(sign string, params *qf.ChargeParams) (*qf.Charge, error) {
	return getC().New(sign, params)
}

// New creates a new charge.
func (c Client) New(sign string, params *qf.ChargeParams) (*qf.Charge, error) {
	path := "/v1/payment"
	charge := &qf.Charge{}

	err := c.B.Call("POST", path, c.AppCode, sign, nil, params, charge)
	return charge, err
}

func getC() Client {
	return Client{qf.GetBackend(qf.APIBackend), qf.AppCode}
}
