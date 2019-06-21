package wxpay

import "testing"

func TestClient_UnifiedOrder(t *testing.T) {
	client := NewClient(NewAccount("wx2421b1c4370ec43b", "10000100", "xxx", false))

	params := make(Params)
	params.SetString("body", "test").
		SetString("out_trade_no", "58867657575757").
		SetInt64("total_fee", 1).
		SetString("spbill_create_ip", "127.0.0.1").
		SetString("notify_url", "http://notify.objcoding.com/notify").
		SetString("trade_type", "APP")

	t.Log(client.UnifiedOrder(params))
}
