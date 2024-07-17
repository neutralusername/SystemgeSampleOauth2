package appWebsocketHTTP

import (
	"Systemge/Oauth2"
)

type AppWebsocketHTTP struct {
	oauth2Server *Oauth2.Server
}

func New(oauth2Server *Oauth2.Server) *AppWebsocketHTTP {
	return &AppWebsocketHTTP{
		oauth2Server: oauth2Server,
	}
}
