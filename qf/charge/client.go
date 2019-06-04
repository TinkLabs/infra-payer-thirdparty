package charge

import (
	"net/http"
	"net/url"

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

	uv := &url.Values{}

	if params.Txamt != "" {
		uv.Set("txamt", params.Txamt)
	}
	if params.Txcurrcd != "" {
		uv.Set("txcurrcd", params.Txcurrcd)
	}
	if params.PayType != "" {
		uv.Set("pay_type", string(params.PayType))
	}
	if params.OutTradeNo != "" {
		uv.Set("out_trade_no", params.OutTradeNo)
	}

	if params.Txdtm != "" {
		uv.Set("txdtm", params.Txdtm)
	}
	if params.GoodsName != "" {
		uv.Set("goods_name", params.GoodsName)
	}

	if params.ExpiredTime > 0 {
		uv.Set("expired_time", params.ExpiredTime)
	}

	if params.Mchid != "" {
		uv.Set("mchid", params.Mchid)
	}
	if params.PayTag != "" {
		uv.Set("pay_tag", params.PayTag)
	}

	if params.Txzone != "" {
		uv.Set("txzone", params.Txzone)
	}
	if params.ReturnUrl != "" {
		uv.Set("return_url", params.ReturnUrl)
	}
	if params.LimitPay != "" {
		uv.Set("limit_pay", params.LimitPay)
	}
	if params.Udid != "" {
		uv.Set("udid", params.Udid)
	}

	err := c.B.Call(http.MethodPost, path, c.AppCode, sign, uv, nil, charge)
	return charge, err
}

func getC() Client {
	return Client{qf.GetBackend(qf.APIBackend), qf.AppCode}
}
