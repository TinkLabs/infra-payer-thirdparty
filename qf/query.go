package qf

// QueryParams is the set of parameters that can be used when creating a query.
type QueryParams struct {
	Mchid      string      `json:"mchid"`
	Syssn      string      `json:"syssn"`
	OutTradeNo string      `json:"out_trade_no"`
	PayType    PaymentType `json:"pay_type"`
	Respcd     PaymentCode `json:"respcd"`
	StartTime  string      `json:"start_time"`
	EndTime    string      `json:"end_time"`
	Txzone     string      `json:"txzone"`
	Page       string      `json:"page"`
	PageSize   string      `json:"page_size"`
}

// Query is the resource representing a qf QueryList.
type Query struct {
	Page     string      `json:"page"`
	Resperr  string      `json:"resperr"`
	PageSize string      `json:"page_size"`
	Respcd   PaymentCode `json:"respcd"`
	Data     []MetaData  `json:"data"`
}

type MetaData struct {
	Syssn      string      `json:"syssn"`
	OutTradeNo string      `json:"out_trade_no"`
	PayType    PaymentType `json:"pay_type"`
	OrderType  string      `json:"order_type"`
	Txdtm      string      `json:"txdtm"`
	Txamt      string      `json:"txamt"`
	Sysdtm     string      `json:"sysdtm"`
	Cancel     string      `json:"cancel"`
	Respcd     PaymentCode `json:"respcd"`
	Errmsg     string      `json:"errmsg"`
}
