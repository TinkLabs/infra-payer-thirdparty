package qf

// RefundParams is the set of parameters that can be used when refunding a charge.
type RefundParams struct {
	Syssn      string `json:"syssn"`
	OutTradeNo string `json:"out_trade_no"`
	Txamt      string `json:"txamt"`
	Txdtm      string `json:"txdtm"`
	Mchid      string `json:"mchid,omitempty"`
	Txzone     string `json:"txzone,omitempty"`
	ReturnUrl  string `json:"return_url,omitempty"`
	Udid       string `json:"udid,omitempty"`
}

// Refund is the resource representing a qf refund.
type Refund struct {
	OrigSyssn  string      `json:"orig_syssn"`
	Sysdtm     string      `json:"sysdtm"`
	Cardcd     string      `json:"cardcd"`
	Txdtm      string      `json:"txdtm"`
	Resperr    string      `json:"resperr"`
	Txcurrcd   string      `json:"txcurrcd"`
	Txamt      string      `json:"txamt"`
	Respmsg    string      `json:"respmsg"`
	OutTradeNo string      `json:"out_trade_no"`
	Syssn      string      `json:"syssn"`
	Respcd     PaymentCode `json:"respcd"`
}
