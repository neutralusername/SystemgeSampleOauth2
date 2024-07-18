package main

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Oauth2"
	"Systemge/Utilities"
	"SystemgeSamplePingPong/appWebsocketHTTP"
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	oauth2Server, err := Oauth2.New(Config.Oauth2{
		Randomizer:                 Utilities.NewRandomizer(Utilities.GetSystemTime()),
		Oauth2State:                Utilities.RandomString(16, Utilities.ALPHA_NUMERIC),
		SessionLifetimeMs:          15000,
		Port:                       8081,
		AuthPath:                   "/",
		AuthCallbackPath:           "/callback",
		SucessUrlCallbackRedirect:  "http://localhost:8080",
		FailureUrlCallbackRedirect: "http://chatgpt.com",
		OAuth2Config: &oauth2.Config{
			ClientID:     "1261641608886222908",
			ClientSecret: "xD",
			RedirectURL:  "http://localhost:8081/callback",
			Scopes:       []string{"identify"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
		TokenHandler: func(oauth2Config *oauth2.Config, token *oauth2.Token) (string, map[string]interface{}, error) {
			client := oauth2Config.Client(context.Background(), token)
			resp, err := client.Get("https://discord.com/api/users/@me")
			if err != nil {
				return "", nil, Error.New("failed getting user", err)
			}
			defer resp.Body.Close()

			var discordAuthData map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&discordAuthData); err != nil {
				return "", nil, Error.New("failed decoding user", err)
			}
			if discordAuthData["username"] == nil {
				return "", nil, Error.New("failed getting session identity", nil)
			}
			return discordAuthData["username"].(string), discordAuthData, nil
		},
	})
	if err != nil {
		panic(err)
	}
	Module.StartCommandLineInterface(Module.NewMultiModule(true,
		Node.New(Config.Node{
			Name:   "nodeOauth2",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, oauth2Server),
		Node.New(Config.Node{
			Name:   "nodeWebsocketHTTP",
			Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
		}, appWebsocketHTTP.New(oauth2Server)),
	))
}
