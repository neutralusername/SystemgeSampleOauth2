package main

import (
	"SystemgeSamplePingPong/appWebsocketHTTP"
	"context"
	"encoding/json"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Oauth2"
	"github.com/neutralusername/Systemge/Tools"

	"golang.org/x/oauth2"
)

const LOGGER_PATH = "logs.log"

var gmailConfig = &Config.Oauth2{
	Oauth2State:       Tools.RandomString(16, Tools.ALPHA_NUMERIC),
	SessionLifetimeMs: 15000,
	Server: &Config.TcpServer{
		Port:        8081,
		TlsCertPath: "MyCertificate.crt",
		TlsKeyPath:  "MyKey.key",
		Blacklist:   []string{},
		Whitelist:   []string{},
	},
	AuthPath:                   "/",
	AuthCallbackPath:           "/callback",
	CallbackSuccessRedirectUrl: "https://localhost:8080",
	CallbackFailureRedirectUrl: "https://chatgpt.com",
	OAuth2Config: &oauth2.Config{
		ClientID:     "489235287049-jdbort0h24p9pfiupqpu8616dvgslq2t.apps.googleusercontent.com", // replace with your own
		ClientSecret: Helpers.GetFileContent("gmailClientSecret.txt"),                            // replace with your own
		RedirectURL:  "https://localhost:8081/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	},
	TokenHandler: func(oauth2Config *oauth2.Config, token *oauth2.Token) (string, map[string]interface{}, error) {
		client := oauth2Config.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			return "", nil, Error.New("failed getting user", err)
		}
		defer resp.Body.Close()

		var googleAuthData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&googleAuthData); err != nil {
			return "", nil, Error.New("failed decoding user", err)
		}
		if googleAuthData["email"] == nil {
			return "", nil, Error.New("failed getting session identity", nil)
		}
		return googleAuthData["email"].(string), googleAuthData, nil
	},
}

var discordConfig = &Config.Oauth2{
	Oauth2State:       Tools.RandomString(16, Tools.ALPHA_NUMERIC),
	SessionLifetimeMs: 15000,
	Server: &Config.TcpServer{
		Port:        8081,
		TlsCertPath: "MyCertificate.crt",
		TlsKeyPath:  "MyKey.key",
		Blacklist:   []string{},
		Whitelist:   []string{},
	},
	AuthPath:                   "/",
	AuthCallbackPath:           "/callback",
	CallbackSuccessRedirectUrl: "https://localhost:8080",
	CallbackFailureRedirectUrl: "https://chatgpt.com",
	OAuth2Config: &oauth2.Config{
		ClientID:     "1261641608886222908",                             // replace with your own
		ClientSecret: Helpers.GetFileContent("discordClientSecret.txt"), // replace with your own
		RedirectURL:  "https://localhost:8081/callback",
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
}

func main() {
	loggerQueue := Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	oauth2Server, err := Oauth2.New(discordConfig)
	if err != nil {
		panic(err)
	}
	Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
	}, Dashboard.New(&Config.Dashboard{
		Server: &Config.TcpServer{
			Port: 8082,
		},
		NodeStatusIntervalMs:           1000,
		NodeSystemgeCounterIntervalMs:  1000,
		NodeWebsocketCounterIntervalMs: 1000,
		NodeBrokerCounterIntervalMs:    1000,
		NodeResolverCounterIntervalMs:  1000,
		HeapUpdateIntervalMs:           1000,
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		AutoStart:                      true,
		AddDashboardToDashboard:        true,
	},
		Node.New(&Config.Node{
			Name:           "nodeOauth2",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeOauth2\"] ", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeOauth2\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeOauth2\"] ", loggerQueue),
		}, oauth2Server),
		Node.New(&Config.Node{
			Name:           "nodeWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger:     Tools.NewLogger("[Info \"nodeWebsocketHTTP\"] ", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeWebsocketHTTP\"] ", loggerQueue),
		}, appWebsocketHTTP.New(oauth2Server)),
	)).StartBlocking()
}
