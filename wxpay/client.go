package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	bodyType           = "application/xml; charset=utf-8"
	defaultHTTPTimeout = 6 * time.Second
)

type Client struct {
	account     *Account
	signType    string
	httpTimeout time.Duration
}

func NewClient(account *Account) *Client {
	return &Client{
		account:     account,
		signType:    MD5,
		httpTimeout: defaultHTTPTimeout,
	}
}

func (c *Client) SetHttpTimeoutMs(t time.Duration) {
	c.httpTimeout = t
}

func (c *Client) SetSignType(signType string) {
	c.signType = signType
}

func (c *Client) SetAccount(account *Account) {
	c.account = account
}

func (c *Client) GetAccount() *Account {
	return c.account
}

func (c *Client) fillRequestData(params Params) Params {
	params["appid"] = c.account.appID
	params["mch_id"] = c.account.mchID
	params["nonce_str"] = nonceStr()
	params["sign_type"] = c.signType
	if c.account.isSandbox {
		key, err := c.GetSignKey()
		if err != nil {
			log.Printf("Sign Error", err)
		}
		params["sign"] = c.Sign(params, key)
	} else {
		params["sign"] = c.Sign(params, c.account.apiKey)
	}

	return params
}

// https no cert post
func (c *Client) postWithoutCert(url string, params Params) (string, error) {
	h := &http.Client{Timeout: c.httpTimeout}
	p := c.fillRequestData(params)

	log.Printf("[Debug] Print all UnifiedOrder params: %v\n", p)

	response, err := h.Post(url, bodyType, strings.NewReader(MapToXml(p)))
	if err != nil {
		log.Printf("Post to %s failed: %s\n", url, err.Error())
		return "", errors.New("http.post failed.")
	}

	defer response.Body.Close()

	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Read response failed: %s\n", err.Error())
		return "", errors.New("Read http response failed.")
	}

	return string(res), nil
}

// https need cert post
func (c *Client) postWithCert(url string, params Params) (string, error) {
	if c.account.apiclientCert == nil || c.account.apiclientKey == nil {
		return "", errors.New("Empty cert or key.")
	}

	trans := WithCertBytes(c.account.apiclientCert, c.account.apiclientKey)
	if trans == nil {
		return "", errors.New("Parses PEM key pair failed.")
	}

	h := &http.Client{Timeout: c.httpTimeout, Transport: trans}

	p := c.fillRequestData(params)
	log.Printf("[Debug] Print all Refund params: %v\n", p)

	response, err := h.Post(url, bodyType, strings.NewReader(MapToXml(p)))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (c *Client) generateSignedXml(params Params) string {
	sign := c.Sign(params, c.account.apiKey)
	params.SetString(Sign, sign)
	return MapToXml(params)
}

func (c *Client) ValidSign(params Params) bool {
	if !params.ContainsKey(Sign) {
		return false
	}
	return params.GetString(Sign) == c.Sign(params, c.account.apiKey)
}

func (c *Client) Sign(params Params, apiKey string) string {
	var keys = make([]string, 0, len(params))

	for k := range params {
		if k != "sign" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.GetString(k)) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params.GetString(k))
			buf.WriteString(`&`)
		}
	}

	buf.WriteString(`key=`)
	buf.WriteString(apiKey)

	var (
		dataMd5    [16]byte
		dataSha256 []byte
		str        string
	)

	switch c.signType {
	case MD5:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:])
	case HMACSHA256:
		h := hmac.New(sha256.New, []byte(c.account.apiKey))
		h.Write(buf.Bytes())
		dataSha256 = h.Sum(nil)
		str = hex.EncodeToString(dataSha256[:])
	}

	return strings.ToUpper(str)
}

func (c *Client) processResponseXml(xmlStr string) (Params, error) {
	var returnCode, resultCode string

	log.Printf("[Debug] Print wxpay response params: %v\n", xmlStr)
	params := XmlToMap(xmlStr)

	log.Printf("[Debug] Print XmlToMap: %v\n", params)

	if !params.ContainsKey("return_code") {
		return nil, errors.New("No return_code in XML.")
	}

	returnCode = params.GetString("return_code")

	if returnCode != Success {
		return nil, errors.New("return_code failed.")
	}

	if !params.ContainsKey("result_code") {
		return nil, errors.New("No result_code in XML.")
	}

	resultCode = params.GetString("result_code")
	if resultCode != Success {
		return nil, errors.New("result_code failed.")
	}

	if !c.account.isSandbox && !c.ValidSign(params) {
		return nil, errors.New("Invalid sign value in XML")
	}

	return params, nil
}

func (c *Client) UnifiedOrder(params Params) (Params, error) {
	var url string

	if c.account.isSandbox {
		url = SandboxUnifiedOrderUrl
	} else {
		url = UnifiedOrderUrl
	}

	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) MicroPay(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxMicroPayUrl
	} else {
		url = MicroPayUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) Refund(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxRefundUrl
	} else {
		url = RefundUrl
	}
	xmlStr, err := c.postWithCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) OrderQuery(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxOrderQueryUrl
	} else {
		url = OrderQueryUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) RefundQuery(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxRefundQueryUrl
	} else {
		url = RefundQueryUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) Reverse(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxReverseUrl
	} else {
		url = ReverseUrl
	}
	xmlStr, err := c.postWithCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) CloseOrder(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxCloseOrderUrl
	} else {
		url = CloseOrderUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) DownloadBill(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxDownloadBillUrl
	} else {
		url = DownloadBillUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)

	p := make(Params)

	if strings.Index(xmlStr, "<") == 0 {
		p = XmlToMap(xmlStr)
		return p, err
	} else {
		p.SetString("return_code", Success)
		p.SetString("return_msg", "ok")
		p.SetString("data", xmlStr)
		return p, err
	}
}

func (c *Client) DownloadFundFlow(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxDownloadFundFlowUrl
	} else {
		url = DownloadFundFlowUrl
	}
	xmlStr, err := c.postWithCert(url, params)

	p := make(Params)

	if strings.Index(xmlStr, "<") == 0 {
		p = XmlToMap(xmlStr)
		return p, err
	} else {
		p.SetString("return_code", Success)
		p.SetString("return_msg", "ok")
		p.SetString("data", xmlStr)
		return p, err
	}
}

func (c *Client) Report(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxReportUrl
	} else {
		url = ReportUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) ShortUrl(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxShortUrl
	} else {
		url = ShortUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

func (c *Client) AuthCodeToOpenid(params Params) (Params, error) {
	var url string
	if c.account.isSandbox {
		url = SandboxAuthCodeToOpenidUrl
	} else {
		url = AuthCodeToOpenidUrl
	}
	xmlStr, err := c.postWithoutCert(url, params)
	if err != nil {
		return nil, err
	}
	return c.processResponseXml(xmlStr)
}

type GetSignKeyReq struct {
	Mch_id    string `xml:"mch_id"`
	Nonce_str string `xml:"nonce_str"`
	Sign      string `xml:"sign"`
}

type GetSignKeyResp struct {
	Return_code     string `xml:"return_code"`
	Return_msg      string `xml:"return_msg"`
	Mch_id          string `xml:"mch_id"`
	Sandbox_signkey string `xml:"sandbox_signkey"`
}

//获取沙箱测试的api key
func (c *Client) GetSignKey() (key string, err error) {
	var req GetSignKeyReq
	req.Mch_id = c.account.mchID
	req.Nonce_str = nonceStr()
	m := make(Params)
	m.SetString("mch_id", c.account.mchID)
	m.SetString("nonce_str", req.Nonce_str)
	req.Sign = c.Sign(m, c.account.apiKey)

	bytesReq, _err := xml.Marshal(req)
	if _err != nil {
		log.Printf("以xml形式编码错误, 原因: %v", err)
		return
	}

	strReq := string(bytesReq)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	strReq = strings.Replace(strReq, "GetSignKeyReq", "xml", -1)
	bytesReq = []byte(strReq)

	// wxpay的getsignkey接口需要http
	signReq, _err := http.NewRequest("POST", SandboxSignURL, bytes.NewReader(bytesReq))
	if _err != nil {
		log.Printf("获取验收仿真测试系统的API验签密钥错误，error: %s", err)
		return
	}
	signReq.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	signReq.Header.Set("Content-Type", "application/xml;charset=utf-8")

	client := http.Client{}
	resp, err := client.Do(signReq)
	if err != nil {
		log.Printf("获取验收仿真测试系统的API验签密钥错误, 原因:", err)
		respData, _err := ioutil.ReadAll(resp.Body)
		if _err != nil {
			log.Printf("error:", _err)
			err = _err
			return
		}
		log.Printf(string(respData))
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	xmlResp := GetSignKeyResp{}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error:", err)
		return
	}
	err = xml.Unmarshal(respData, &xmlResp)
	if err != nil {
		log.Printf("error:", err)
		return
	}
	//处理return code.
	if xmlResp.Return_code == "FAIL" {
		log.Printf("获取验收仿真测试系统的API验签密钥错误，原因:Code=%v MSG=%v", xmlResp.Return_code, xmlResp.Return_msg)
		return
	}
	key = xmlResp.Sandbox_signkey
	return
}
