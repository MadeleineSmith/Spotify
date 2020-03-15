package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type CallbackHandler struct {
	HTTPClient *http.Client
}

type TokenRequestBody struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (h CallbackHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rawQuery := req.URL.RawQuery
	v, _ := url.ParseQuery(rawQuery)

	// TODO code isn't available if user hasn't accepted the request
	code := v["code"][0]

	URLdata := url.Values{}
	URLdata.Set("grant_type", "authorization_code")
	URLdata.Set("code", code)
	URLdata.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	URLdata.Set("client_id", os.Getenv("CLIENT_ID"))
	URLdata.Set("client_secret", os.Getenv("CLIENT_SECRET"))

	spotifyRequest, _ := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(URLdata.Encode()))
	spotifyRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	spotifyResponse, _ := h.HTTPClient.Do(spotifyRequest)
	spotifyResponseBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)

	var tokenResponse TokenResponse
	json.Unmarshal(spotifyResponseBodyBytes, &tokenResponse)

	http.Redirect(w, req, fmt.Sprintf("%s/createPlaylist/%s", os.Getenv("FE_BASE_URL"), tokenResponse.AccessToken), http.StatusSeeOther)
}
