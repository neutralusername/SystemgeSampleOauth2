package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Message"
	"Systemge/Node"
	"net/http"

	"github.com/gorilla/websocket"
)

func (app *AppWebsocketHTTP) GetWebsocketComponentConfig() *Config.Websocket {
	return &Config.Websocket{
		Pattern: "/ws",
		Server: &Config.TcpServer{
			Port:      8443,
			Blacklist: []string{},
			Whitelist: []string{},
		},
		HandleClientMessagesSequentially: false,
		ClientMessageCooldownMs:          0,
		ClientWatchdogTimeoutMs:          20000,
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		"authAttempt": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			session := app.oauth2Server.GetSession(message.GetPayload())
			if session == nil {
				websocketClient.Send(Message.NewAsync("authFailure", node.GetName(), "session not found").Serialize())
				return nil
			}
			websocketClient.Send(Message.NewAsync("authSuccess", node.GetName(), session.GetIdentity()).Serialize())
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
