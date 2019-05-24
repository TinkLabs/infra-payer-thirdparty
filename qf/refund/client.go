package refund

import (
	"net/http"
	"net/url"

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

	uv := &url.Values{}

	if params.Syssn != "" {
		uv.Set("syssn", params.Syssn)
	}
	if params.OutTradeNo != "" {
		uv.Set("out_trade_no", params.OutTradeNo)
	}

	if params.Txamt != "" {
		uv.Set("txamt", params.Txamt)
	}
	if params.Txdtm != "" {
		uv.Set("txdtm", params.Txdtm)
	}

	if params.Mchid != "" {
		uv.Set("mchid", params.Mchid)
	}
	if params.Txzone != "" {
		uv.Set("txzone", params.Txzone)
	}

	if params.ReturnUrl != "" {
		uv.Set("return_url", params.ReturnUrl)
	}
	if params.Udid != "" {
		uv.Set("udid", params.Udid)
	}

	err := c.B.Call(http.MethodPost, path, c.AppCode, sign, uv, nil, refund)
	return refund, err
}

func getC() Client {
	return Client{qf.GetBackend(qf.APIBackend), qf.AppCode}
}
