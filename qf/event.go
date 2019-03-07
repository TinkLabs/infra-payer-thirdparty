package qf

// Event is the resource representing a Qf event.
// For more details see https://doc.qfapi.com/docs/ZH-CN/#developer/#_6.
type PaymentEvent struct {
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

// (fu*k qf) Actual payload ex:
// {
// 	"sysdtm": "2019-03-06 15:05:21",
// 	"txcurrcd": "CNY",
// 	"orig_out_trade_no": "user1/eccf606c-0d2a-4ff8-ad09-aeb964b4edd5",
// 	"mchid": "R1zQrTdJnn",
// 	"txdtm": "2019-03-06 15:02:09",
// 	"txamt": "1",
// 	"orig_syssn": "20190306000200020078674642",
// 	"out_trade_no": "user1/1907asdfaecc57adfaassdfa15-81/refund",
// 	"syssn": "20190306000300020078700582",
// 	"respcd": "0000",
// 	"notify_type": "refund"
// }
type RefundEvent struct {
	NotifyType     string      `json:"notify_type"`
	Sysdtm         string      `json:"sysdtm"`
	Txcurrcd       string      `json:"txcurrcd"`
	OrigOutTradeNo string      `json:"orig_out_trade_no"`
	Mchid          string      `json:"mchid"`
	Txdtm          string      `json:"txdtm"`
	Txamt          string      `json:"txamt"`
	OrigSyssn      string      `json:"orig_syssn"`
	OutTradeNo     string      `json:"out_trade_no"`
	Syssn          string      `json:"syssn"`
	Respcd         PaymentCode `json:"respcd"`
}
