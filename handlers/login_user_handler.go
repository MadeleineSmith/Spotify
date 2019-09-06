package handlers

import (
	. "Spotify/constants"
	"net/http"
	"net/url"
)

type LoginUserHandler struct {
	HTTPClient *http.Client
}


// app.get('/login', function(req, res) {
//
//  var state = generateRandomString(16);
//  res.cookie(stateKey, state);
//
//  // your application requests authorization
//  var scope = 'user-read-private user-read-email playlist-modify-private';
//
//  var authorizeURL = 'https://accounts.spotify.com/authorize?' +
//      querystring.stringify({
//          response_type: 'code',
//          client_id: client_id,
//          scope: scope,
//          redirect_uri: redirect_uri,
//          state: state
//      });
//
//
//  res.redirect(authorizeURL);
//});


//  spotifyRequest, _ := http.NewRequest(http.MethodGet, url, nil)
//	spotifyRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUTHORIZATION_TOKEN))
//
//	spotifyResponse, _ := h.HTTPClient.Do(spotifyRequest)
//
//	spotifyBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)
//	spotifyResponseBody := new(SpotifyResponse)
//	json.Unmarshal(spotifyBodyBytes, &spotifyResponseBody)
//
//	// TODO - do I want to remove this track from the JSON I return?
//	if len(spotifyResponseBody.Tracks.Items) == 1 {
//		track.URI = spotifyResponseBody.Tracks.Items[0].URI
//	}

func (h LoginUserHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	authorizeUrl, _ := url.Parse("https://accounts.spotify.com/authorize")
	scopes := "user-read-private user-read-email playlist-modify-private"

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", CLIENT_ID)
	params.Add("scope", scopes)
	params.Add("redirect_uri", REDIRECT_URI)
	// TODO - add `state` query param

	authorizeUrl.RawQuery = params.Encode()

	http.Redirect(w, req, authorizeUrl.String(), http.StatusFound)

	//spotifyRequest, _ := http.NewRequest(http.MethodGet, authorizeUrl.String(), nil)
	//resp, err := h.HTTPClient.Do(spotifyRequest)

	////println(resp.Body)
	//bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//bodyString := string(bodyBytes)
	////print raw response body for debugging purposes
	//fmt.Println("\n\n", bodyString, "\n\n")
	//
	//if err != nil {
	//	println(err.Error())
	//}


	// TODO - still struggling with CORS?
	// try figure out why this isn't working

}