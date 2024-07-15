package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/TcpServer"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		"authAttempt": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			session := app.oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("authFailure", node.GetName(), "session not found").Serialize())
				return nil
			}
			username, ok := session.Get("username")
			if !ok {
				websocketClient.Send(Message.NewAsync("authFailure", node.GetName(), "username not found").Serialize())
				return nil
			}
			websocketClient.Send(Message.NewAsync("authSuccess", node.GetName(), username.(string)).Serialize())
			return nil
		},
		"logoutAttempt": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			session := app.oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("logoutFailure", node.GetName(), "session not found").Serialize())
				return nil
			}
			app.oauth2Server.Expire(session)
			websocketClient.Send(Message.NewAsync("logoutSuccess", node.GetName(), "").Serialize())
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {

}

func (app *AppWebsocketHTTP) GetWebsocketComponentConfig() Config.Websocket {
	return Config.Websocket{
		Pattern:                          "/ws",
		Server:                           TcpServer.New(8443, "", ""),
		HandleClientMessagesSequentially: false,
		ClientMessageCooldownMs:          0,
		ClientWatchdogTimeoutMs:          20000,
	}
}
