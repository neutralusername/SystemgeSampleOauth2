package appWebsocketHTTP

import (
	"Systemge/Node"
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

func (app *AppWebsocketHTTP) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{}
}
