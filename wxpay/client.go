package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	bodyType           = "application/xml; charset=utf-8"
	defaultHTTPTimeout = 10 * time.Second
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

func (c *Client) fillRequestData(params Params) Params {
	params["appid"] = c.account.appID
	params["mch_id"] = c.account.mchID
	params["nonce_str"] = nonceStr()
	params["sign_type"] = c.signType
	params["sign"] = c.Sign(params)

	return params
}

// https no cert post
func (c *Client) postWithoutCert(url string, params Params) (string, error) {
	h := &http.Client{Timeout: c.httpTimeout}
	p := c.fillRequestData(params)

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

// https need cert post
func (c *Client) postWithCert(url string, params Params) (string, error) {
	if c.account.certData == nil {
		return "", errors.New("Empty cert data.")
	}

	cert := pkcs12ToPem(c.account.certData, c.account.mchID)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	transport := &http.Transport{
		TLSClientConfig:    config,
		DisableCompression: true,
	}
	h := &http.Client{Timeout: c.httpTimeout, Transport: transport}
	p := c.fillRequestData(params)
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
	sign := c.Sign(params)
	params.SetString(Sign, sign)
	return MapToXml(params)
}

func (c *Client) ValidSign(params Params) bool {
	if !params.ContainsKey(Sign) {
		return false
	}
	return params.GetString(Sign) == c.Sign(params)
}

func (c *Client) Sign(params Params) string {
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
	buf.WriteString(c.account.apiKey)

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
	params := XmlToMap(xmlStr)

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

	if !c.ValidSign(params) {
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
