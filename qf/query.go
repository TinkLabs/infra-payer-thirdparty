package qf

// QueryParams is the set of parameters that can be used when creating a query.
type QueryParams struct {
	Txamt      string `json:"txamt,omitempty"`
	Txcurrcd   string `json:"txcurrcd,omitempty"`
	PayType    string `json:"pay_type,omitempty"`
	OutTradeNo string `json:"out_trade_no,omitempty"`
	Txdtm      string `json:"txdtm,omitempty"`
	GoodsName  string `json:"goods_name,omitempty"`
	Mchid      string `json:"mchid,omitempty"`
	PayTag     string `json:"pay_tag,omitempty"`
	Txzone     string `json:"txzone,omitempty"`
	ReturnUrl  string `json:"return_url,omitempty"`
	LimitPay   string `json:"limit_pay,omitempty"`
	Udid       string `json:"udid,omitempty"`
}

// Query is the resource representing a qf QueryList.
type Query struct {
	PayType    string `json:"pay_type"`
	Sysdtm     string `json:"sysdtm"`
	Txdtm      string `json:"txdtm"`
	Resperr    string `json:"resperr"`
	Txamt      string `json:"txamt"`
	Respmsg    string `json:"respmsg"`
	OutTradeNo string `json:"out_trade_no"`
	Syssn      string `json:"syssn"`
	Qrcode     string `json:"qrcode"`
	Respcd     string `json:"respcd"`
}
