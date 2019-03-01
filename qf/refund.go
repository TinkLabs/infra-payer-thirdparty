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
	Syssn     string `json:"syssn"`
	OrigSyssn string `json:"orig_syssn"`
	Txamt     string `json:"txamt"`
	Txdtm     string `json:"txdtm"`
	Sysdtm    string `json:"sysdtm"`
}
