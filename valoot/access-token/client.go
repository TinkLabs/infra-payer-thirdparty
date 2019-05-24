package access_token

import (
	"github.com/TinkLabs/payer-thirdparty/valoot"
	"net/http"
)

type Client struct {
	B valoot.Backend
}

func GetAccessToken(params *valoot.AccessTokenParams) (*valoot.AccessTokenResp, error) {
	return getC().GetAccessToken(params)
}

func RefreshAccessToken(params *valoot.AccessTokenParams) (*valoot.AccessTokenResp, error) {
	return getC().RefreshAccessToken(params)
}

func (c Client) GetAccessToken(params *valoot.AccessTokenParams) (resp *valoot.AccessTokenResp, err error) {
	params.GrantType = "password"

	path := "/v1/oauth/token"
	err = c.B.Call(http.MethodPost, path, "", nil, params, &resp)
	return resp, err
}

func (c Client) RefreshAccessToken(params *valoot.AccessTokenParams) (resp *valoot.AccessTokenResp, err error) {
	params.GrantType = "refresh_token"

	path := "/v1/oauth/token"
	err = c.B.Call(http.MethodPost, path, "", nil, params, &resp)
	return resp, err
}

func getC() Client {
	return Client{valoot.GetBackend(valoot.APIBackend)}
}
