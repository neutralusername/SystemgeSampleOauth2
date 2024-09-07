package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Oauth2Server"
	"github.com/neutralusername/Systemge/Tools"
	"github.com/neutralusername/Systemge/WebsocketServer"

	"golang.org/x/oauth2"
)

const LOGGER_PATH = "logs.log"

var gmailConfig = &Config.Oauth2{
	TcpServerConfig: &Config.TcpServer{
		Port:        8082,
		TlsCertPath: "MyCertificate.crt",
		TlsKeyPath:  "MyKey.key",
		Blacklist:   []string{},
		Whitelist:   []string{},
	},
	Oauth2State:                Tools.GenerateRandomString(16, Tools.ALPHA_NUMERIC),
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
	TcpServerConfig: &Config.TcpServer{
		Port:        8082,
		TlsCertPath: "MyCertificate.crt",
		TlsKeyPath:  "MyKey.key",
		Blacklist:   []string{},
		Whitelist:   []string{},
	},
	Oauth2State:                Tools.GenerateRandomString(16, Tools.ALPHA_NUMERIC),
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
	oauth2Server := Oauth2Server.New("oauth2Server", discordConfig)
	oauth2Server.Start()
	websocketServer := WebsocketServer.New("websocketServer",
		&Config.WebsocketServer{
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
			MailerConfig:      nil,
			Pattern:           "/ws",
			TcpServerConfig: &Config.TcpServer{
				Port:        8443,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ClientRateLimiterBytes:           nil,
			ClientRateLimiterMessages:        nil,
			IncomingMessageByteLimit:         0,
			HandleClientMessagesSequentially: false,
			ClientWatchdogTimeoutMs:          60000,
			Upgrader: &websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		}, getWebsocketMessageHandlers(oauth2Server), nil, nil,
	)
	go func() {
		err := websocketServer.Start()
		if err != nil {
			panic(err)
		}
	}()
	httpServer := HTTPServer.New("httpServer",
		&Config.HTTPServer{
			TcpServerConfig: &Config.TcpServer{
				Port:        8080,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
		}, GgtHTTPMessageHandlers(),
	)
	httpServer.Start()
	time.Sleep(1000 * time.Hour)
}

func GgtHTTPMessageHandlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/": HTTPServer.SendDirectory("../frontend"),
	}
}

func getWebsocketMessageHandlers(oauth2Server *Oauth2Server.Server) WebsocketServer.MessageHandlers {
	return map[string]WebsocketServer.MessageHandler{
		"authAttempt": func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
			session := oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("authFailure", "session not found").Serialize())
				return nil
			}
			websocketClient.Send(Message.NewAsync("authSuccess", session.GetIdentity()).Serialize())
			return nil
		},
		"logoutAttempt": func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
			session := oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("logoutFailure", "session not found").Serialize())
				return nil
			}
			oauth2Server.Expire(session)
			websocketClient.Send(Message.NewAsync("logoutSuccess", "").Serialize())
			return nil
		},
	}
}
