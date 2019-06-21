package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	appID     string
	mchID     string
	apiKey    string
	certData  []byte
	isSandbox bool
}

func NewAccount(appID string, mchID string, apiKey string, isSanbox bool) *Account {
	return &Account{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		isSandbox: isSanbox,
	}
}

func (a *Account) SetCertData(certPath string) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Println("Failed to read cert file.")
		return
	}
	a.certData = certData
}
