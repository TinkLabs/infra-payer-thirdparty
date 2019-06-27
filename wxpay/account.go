package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	appID         string
	mchID         string
	apiKey        string
	apiclientCert []byte
	apiclientKey  []byte
	isSandbox     bool
}

func NewAccount(appID string, mchID string, apiKey string, isSanbox bool) *Account {
	return &Account{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		isSandbox: isSanbox,
	}
}

func (a *Account) SetApplicationCert(certPath string) {
	appCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Println("Failed to read cert file.")
		return
	}

	a.apiclientCert = appCert
}

func (a *Account) SetApplicationKey(keyPath string) {
	appkey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Println("Failed to read key file.")
		return
	}

	a.apiclientKey = appkey
}
