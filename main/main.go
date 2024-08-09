package main

import (
	"SystemgeSampleOauth2/appWebsocketHTTP"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
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
	NodeConfig: &Config.Node{
		Name:              "nodeOauth2",
		RandomizerSeed:    Tools.GetSystemTime(),
		InfoLoggerPath:    LOGGER_PATH,
		WarningLoggerPath: LOGGER_PATH,
		ErrorLoggerPath:   LOGGER_PATH,
	},
	ServerConfig: &Config.TcpServer{
		Port:        8082,
		TlsCertPath: "MyCertificate.crt",
		TlsKeyPath:  "MyKey.key",
		Blacklist:   []string{},
		Whitelist:   []string{},
	},
	Oauth2State:                Tools.RandomString(16, Tools.ALPHA_NUMERIC),
	SessionLifetimeMs:          15000,
	AuthPath:                   "/",
	AuthCallbackPath:           "/callback",
	CallbackSuccessRedirectUrl: "https://localhost:8080",
	CallbackFailureRedirectUrl: "https://chatgpt.com",
	OAuth2Config: &oauth2.Config{
		ClientID:     "489235287049-jdbort0h24p9pfiupqpu8616dvgslq2t.apps.googleusercontent.com", // replace with your own
		ClientSecret: Helpers.GetFileContent("gmailClientSecret.txt"),                            // replace with your own
		RedirectURL:  "https://localhost:8082/callback",
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
	NodeConfig: &Config.Node{
		Name:              "nodeOauth2",
		RandomizerSeed:    Tools.GetSystemTime(),
		InfoLoggerPath:    LOGGER_PATH,
		WarningLoggerPath: LOGGER_PATH,
		ErrorLoggerPath:   LOGGER_PATH,
	},
	ServerConfig: &Config.TcpServer{
		Port:        8082,
		TlsCertPath: "MyCertificate.crt",
		TlsKeyPath:  "MyKey.key",
		Blacklist:   []string{},
		Whitelist:   []string{},
	},
	Oauth2State:                Tools.RandomString(16, Tools.ALPHA_NUMERIC),
	SessionLifetimeMs:          15000,
	AuthPath:                   "/",
	AuthCallbackPath:           "/callback",
	CallbackSuccessRedirectUrl: "https://localhost:8080",
	CallbackFailureRedirectUrl: "https://chatgpt.com",
	OAuth2Config: &oauth2.Config{
		ClientID:     "1261641608886222908",                             // replace with your own
		ClientSecret: Helpers.GetFileContent("discordClientSecret.txt"), // replace with your own
		RedirectURL:  "https://localhost:8082/callback",
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
	Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	oauth2Node, err := Oauth2.New(discordConfig)
	if err != nil {
		panic(err)
	}
	Dashboard.New(&Config.Dashboard{
		NodeConfig: &Config.Node{
			Name:           "dashboard",
			RandomizerSeed: Tools.GetSystemTime(),
		},
		ServerConfig: &Config.TcpServer{
			Port: 8081,
		},
		NodeStatusIntervalMs:           1000,
		NodeSystemgeCounterIntervalMs:  1000,
		NodeWebsocketCounterIntervalMs: 1000,
		HeapUpdateIntervalMs:           1000,
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		GoroutineUpdateIntervalMs:      1000,
		AutoStart:                      true,
		AddDashboardToDashboard:        true,
	},
		oauth2Node,
		Node.New(&Config.NewNode{
			HttpConfig: &Config.HTTP{
				ServerConfig: &Config.TcpServer{
					Port:        8080,
					TlsCertPath: "MyCertificate.crt",
					TlsKeyPath:  "MyKey.key",
				},
			},
			WebsocketConfig: &Config.Websocket{
				Pattern: "/ws",
				ServerConfig: &Config.TcpServer{
					Port:        8443,
					TlsCertPath: "MyCertificate.crt",
					TlsKeyPath:  "MyKey.key",
					Blacklist:   []string{},
					Whitelist:   []string{},
				},
				HandleClientMessagesSequentially: false,
				ClientWatchdogTimeoutMs:          20000,
				Upgrader: &websocket.Upgrader{
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
			NodeConfig: &Config.Node{
				Name:              "nodeWebsocketHTTP",
				RandomizerSeed:    Tools.GetSystemTime(),
				InfoLoggerPath:    LOGGER_PATH,
				WarningLoggerPath: LOGGER_PATH,
				ErrorLoggerPath:   LOGGER_PATH,
			},
		}, appWebsocketHTTP.New(oauth2Node.GetApplication().(*Oauth2.App))),
	).StartBlocking()
}
