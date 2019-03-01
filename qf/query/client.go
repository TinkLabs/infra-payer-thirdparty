package query

import (
	"github.com/TinkLabs/payer-thirdparty/qf"
)

type Client struct {
	B       qf.Backend
	AppCode string
}

func Get(sign string, params *qf.QueryParams) (*qf.Query, error) {
	return getC().Get(sign, params)
}

func (c Client) Get(sign string, params *qf.QueryParams) (*qf.Query, error) {
	path := "/v1/query"
	query := &qf.Query{}

	err := c.B.Call("POST", path, c.AppCode, sign, nil, params, query)
	return query, err
}

func getC() Client {
	return Client{qf.GetBackend(qf.APIBackend), qf.AppCode}
}
