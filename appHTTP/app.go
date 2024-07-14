package appHTTP

import (
	"Systemge/Config"
	"Systemge/Node"
	"Systemge/Oauth2"
)

type AppHTTP struct {
	oauth2Server *Oauth2.Server
}

func New(oauth2Server *Oauth2.Server) *AppHTTP {
	return &AppHTTP{
		oauth2Server: oauth2Server,
	}
}

func (app *AppHTTP) OnStart(node *Node.Node) error {
	return nil
}

func (app *AppHTTP) OnStop(node *Node.Node) error {
	return nil
}

func (app *AppHTTP) GetApplicationConfig() Config.Application {
	return Config.Application{
		HandleMessagesSequentially: false,
	}
}
