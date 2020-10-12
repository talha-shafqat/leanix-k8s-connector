package leanix

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// AuthResponse struct
type AuthResponse struct {
	Scope       string `json:"scope"`
	Expired     bool   `json:"expired"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Authenticate uses token to authenticate against integration api and reponse with access_token
func Authenticate(fqdn string, token string) (string, error) {
	body := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", "https://"+fqdn+"/services/mtm/v1/oauth2/token", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("apitoken", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	authResponse := AuthResponse{}
	json.Unmarshal(responseData, &authResponse)
	return authResponse.AccessToken, nil
}

func Upload(fqdn string, accessToken string, ldif []byte) {

}

func StartRun() {

}
