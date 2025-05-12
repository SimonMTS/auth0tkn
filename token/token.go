package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"s14.nl/auth0tkn/profile"
)

type Token struct {
	Prefix,
	Token string
	ValidUntil int
}

func Get(p profile.Profile) (Token, error) {
	data := url.Values{}

	data.Set("client_id", p.Tenant.ClientId)
	data.Set("client_secret", p.Tenant.ClientSecret)
	data.Set("audience", p.Tenant.Audience)
	data.Set("scope", "openid email profile")

	data.Set("grant_type", "password")
	data.Set("username", p.Username)
	data.Set("password", p.Password)

	r, err := http.NewRequest(
		http.MethodPost,
		p.Tenant.Url+"/oauth/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return Token{}, err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return Token{}, err
	}
	if resp.StatusCode != 200 {
		return Token{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		Scopes       string `json:"scopes"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Token{}, err
	}
	if result.AccessToken == "" {
		return Token{}, fmt.Errorf("empty access token: response body: %s", body)
	}

	return Token{
		Prefix:     result.TokenType,
		Token:      result.AccessToken,
		ValidUntil: int(time.Now().Unix()) + result.ExpiresIn,
	}, nil
}
