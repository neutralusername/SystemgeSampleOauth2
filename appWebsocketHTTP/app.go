package appWebsocketHTTP

import (
	"Systemge/Config"
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

func (app *AppWebsocketHTTP) OnStart(node *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(node *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) GetApplicationConfig() Config.Application {
	return Config.Application{
		HandleMessagesSequentially: false,
	}
}
