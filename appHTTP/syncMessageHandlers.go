package appHTTP

import (
	"Systemge/Node"
)

func (app *AppHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}
