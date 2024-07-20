package appWebsocketHTTP

import (
	"Systemge/Oauth2"
)

type AppWebsocketHTTP struct {
	oauth2Server *Oauth2.App
}

func New(oauth2Server *Oauth2.App) *AppWebsocketHTTP {
	return &AppWebsocketHTTP{
		oauth2Server: oauth2Server,
	}
}
