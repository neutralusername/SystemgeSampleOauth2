package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/TcpServer"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		"auth": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			session := app.oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("error", node.GetName(), "Invalid session").Serialize())
				return nil
			}
			username, ok := session.Get("username")
			if !ok {
				websocketClient.Send(Message.NewAsync("error", node.GetName(), "Invalid session").Serialize())
				return nil
			}
			websocketClient.Send(Message.NewAsync("username", node.GetName(), username.(string)).Serialize())
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	websocketClient.Send(Message.NewAsync("auth", node.GetName(), "").Serialize())
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
