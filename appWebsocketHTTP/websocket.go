package appWebsocketHTTP

import (
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		"authAttempt": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			session := app.oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("authFailure", "session not found").Serialize())
				return nil
			}
			websocketClient.Send(Message.NewAsync("authSuccess", session.GetIdentity()).Serialize())
			return nil
		},
		"logoutAttempt": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			session := app.oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("logoutFailure", "session not found").Serialize())
				return nil
			}
			app.oauth2Server.Expire(session)
			websocketClient.Send(Message.NewAsync("logoutSuccess", "").Serialize())
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {

}
