package appHTTP

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Node"
	"Systemge/Oauth2"
	"SystemgeSamplePingPong/topics"
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
	err := node.AsyncMessage(topics.PING, node.GetName(), "ping")
	if err != nil {
		node.GetLogger().Error(Error.New("error sending ping message", err).Error())
	}
	return nil
}

func (app *AppHTTP) OnStop(node *Node.Node) error {
	err := node.AsyncMessage(topics.PING, node.GetName(), "ping")
	if err != nil {
		node.GetLogger().Error(Error.New("error sending ping message", err).Error())
	}
	println("successfully sent ping message to broker but app's node already stopped due to multi-module stop order.")
	return nil
}

func (app *AppHTTP) GetApplicationConfig() Config.Application {
	return Config.Application{
		HandleMessagesSequentially: false,
	}
}
