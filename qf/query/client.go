package query

import (
	"net/http"
	"net/url"

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

	uv := &url.Values{}

	if params.Mchid != "" {
		uv.Set("mchid", params.Mchid)
	}

	if params.Syssn != "" {
		uv.Set("syssn", params.Syssn)
	}

	if params.OutTradeNo != "" {
		uv.Set("out_trade_no", params.OutTradeNo)
	}

	if params.PayType != "" {
		uv.Set("pay_type", string(params.PayType))
	}

	if params.Respcd != "" {
		uv.Set("respcd", string(params.Respcd))
	}

	if params.StartTime != "" {
		uv.Set("start_time", params.StartTime)
	}

	if params.EndTime != "" {
		uv.Set("end_time", params.EndTime)
	}

	if params.Txzone != "" {
		uv.Set("txzone", params.Txzone)
	}

	if params.Page != "" {
		uv.Set("page", params.Page)
	}

	if params.PageSize != "" {
		uv.Set("page_size", params.PageSize)
	}

	err := c.B.Call(http.MethodGet, path, c.AppCode, sign, uv, nil, query)
	return query, err
}

func getC() Client {
	return Client{qf.GetBackend(qf.APIBackend), qf.AppCode}
}
