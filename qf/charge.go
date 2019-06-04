package qf

type PaymentCode string

const (
	PaymentSuccessful PaymentCode = "0000"
	PaymentPendingI   PaymentCode = "1143"
	PaymentPendingII  PaymentCode = "1145"
)

type PaymentType string

const (
	WeChatFrontScan PaymentType = "800201"
	AlipayFrontScan PaymentType = "800101"

	WeChatVersoScan PaymentType = "800208"
	AlipayVersoScan PaymentType = "800108"
)

// ChargeParams is the set of parameters that can be used when creating a charge.
type ChargeParams struct {
	Txamt       string      `json:"txamt"`
	Txcurrcd    string      `json:"txcurrcd"`
	PayType     PaymentType `json:"pay_type"`
	OutTradeNo  string      `json:"out_trade_no"`
	Txdtm       string      `json:"txdtm"`
	GoodsName   string      `json:"goods_name"`
	expiredTime int         `json:"expired_time"`
	Mchid       string      `json:"mchid,omitempty"`
	PayTag      string      `json:"pay_tag,omitempty"`
	Txzone      string      `json:"txzone,omitempty"`
	ReturnUrl   string      `json:"return_url,omitempty"`
	LimitPay    string      `json:"limit_pay,omitempty"`
	Udid        string      `json:"udid,omitempty"`
}

// Charge is the resource representing a qf charge.
type Charge struct {
	PayType    PaymentType `json:"pay_type"`
	Sysdtm     string      `json:"sysdtm"`
	Txdtm      string      `json:"txdtm"`
	Resperr    string      `json:"resperr"`
	Txamt      string      `json:"txamt"`
	Respmsg    string      `json:"respmsg"`
	OutTradeNo string      `json:"out_trade_no"`
	Syssn      string      `json:"syssn"`
	Qrcode     string      `json:"qrcode"`
	Respcd     PaymentCode `json:"respcd"`
}
