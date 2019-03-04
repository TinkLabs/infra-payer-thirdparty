package qf

// Event is the resource representing a Qf event.
// For more details see https://doc.qfapi.com/docs/ZH-CN/#developer/#_6.
type Event struct {
	NotifyType string      `json:"notify_type"`
	Syssn      string      `json:"syssn"`
	OutTradeNo string      `json:"out_trade_no"`
	PayType    PaymentType `json:"pay_type"`
	Txdtm      string      `json:"txdtm"`
	Txamt      string      `json:"txamt"`
	Respcd     PaymentCode `json:"respcd"`
	Sysdtm     string      `json:"sysdtm"`
	Paydtm     string      `json:"paydtm"`
	Cancel     string      `json:"cancel"`
	Cardcd     string      `json:"cardcd"`
	GoodsName  string      `json:"goods_name"`
	Status     string      `json:"status"`
	Txcurrcd   string      `json:"txcurrcd"`
	Mchid      string      `json:"mchid"`
}
