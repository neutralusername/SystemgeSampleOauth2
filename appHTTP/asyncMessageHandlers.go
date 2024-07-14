package appHTTP

import (
	"Systemge/Node"
)

func (app *AppHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}
