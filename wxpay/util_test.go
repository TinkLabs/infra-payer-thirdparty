package wxpay

import (
	"testing"
)

func TestXmlToMap(t *testing.T) {
	// xmlStr := "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg><appid><![CDATA[wx2421b1c4370ec43b]]></appid><mch_id><![CDATA[10000100]]></mch_id><nonce_str><![CDATA[IITRi8Iabbblz1Jc]]></nonce_str><sign><![CDATA[7921E432F65EB8ED0CE9755F0E86D72F]]></sign><result_code><![CDATA[SUCCESS]]></result_code><prepay_id><![CDATA[wx201411101639507cbf6ffd8b0779950874]]></prepay_id><trade_type><![CDATA[APP]]></trade_type></xml>"
	xmlStr := `<xml>
	  <trade_type><![CDATA[APP]]></trade_type>
	  <prepay_id><![CDATA[wx20190506151427404120]]></prepay_id>
	  <nonce_str><![CDATA[1557126867008139000]]></nonce_str>
	  <return_code><![CDATA[SUCCESS]]></return_code>
	  <err_code_des><![CDATA[ok]]></err_code_des>
	  <sign><![CDATA[3C323ECE220CC631240E874F80C3737D]]></sign>
	  <mch_id><![CDATA[1225268501]]></mch_id>
	  <return_msg><![CDATA[OK]]></return_msg>
	  <appid><![CDATA[wxe353414396c4b931]]></appid>
	  <device_info><![CDATA[sandbox]]></device_info>
	  <result_code><![CDATA[SUCCESS]]></result_code>
	  <err_code><![CDATA[SUCCESS]]></err_code>
	</xml>`
	params := XmlToMap(xmlStr)
	if params == nil {
		t.Error(params)
	}
	t.Log(params)
}

func TestMapToXml(t *testing.T) {
	params := map[string]string{"return_msg": "OK", "appid": "wx2421b1c4370ec43b", "mch_id": "10000100"}
	xmlStr := MapToXml(params)
	t.Log(xmlStr)
}

func TestNonceStr(t *testing.T) {
	t.Log(nonceStr())
}
