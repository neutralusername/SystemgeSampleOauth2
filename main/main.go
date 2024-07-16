package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Oauth2"
	"Systemge/Resolver"
	"Systemge/Utilities"
	"SystemgeSamplePingPong/appWebsocketHTTP"
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	randomizer := Utilities.NewRandomizer(Utilities.GetSystemTime())
	oauth2Server, err := (Oauth2.Config{
		Name:                    "discordAuth",
		Randomizer:              randomizer,
		Oauth2State:             randomizer.GenerateRandomString(16, Utilities.ALPHA_NUMERIC),
		SessionLifetimeMs:       15000,
		Port:                    8081,
		AuthPath:                "/",
		AuthCallbackPath:        "/callback",
		SucessCallbackRedirect:  "http://localhost:8080",
		FailureCallbackRedirect: "http://chatgpt.com",
		OAuth2Config: &oauth2.Config{
			ClientID:     "1261641608886222908",
			ClientSecret: "RB9SMRZHm2-JLrkZyzbwj8s-d8S25VTI",
			RedirectURL:  "http://localhost:8081/callback",
			Scopes:       []string{"identify"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
		Logger: Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, nil),
		TokenHandler: func(oauth2Server *Oauth2.Server, token *oauth2.Token) (string, map[string]interface{}, error) {
			client := oauth2Server.GetOauth2Config().Client(context.Background(), token)
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
	}).NewServer()
	if err != nil {
		panic(err)
	}

	err = Resolver.New(Config.ParseResolverConfigFromFile("resolver.systemge")).Start()
	if err != nil {
		panic(err)
	}
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Broker.New(Config.ParseBrokerConfigFromFile("brokerHTTP.systemge")),
		oauth2Server,
		Node.New(Config.ParseNodeConfigFromFile("nodeHTTP.systemge"), appWebsocketHTTP.New(oauth2Server)),
	))
}
