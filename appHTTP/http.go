package appHTTP

import (
	"Systemge/Config"
	"Systemge/Http"
	"Systemge/TcpServer"
	"net/http"
)

func (app *AppHTTP) GetHTTPRequestHandlers() map[string]Http.RequestHandler {
	return map[string]Http.RequestHandler{
		"/": func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
			sessionId := httpRequest.URL.Query().Get("sessionId")
			session := app.oauth2Server.GetSession(sessionId)
			if session == nil {
				responseWriter.Write([]byte("invalid session"))
				return
			}
			username, ok := session.Get("username")
			if !ok {
				responseWriter.Write([]byte("invalid session"))
				return
			}
			responseWriter.Write([]byte("Hello " + username.(string)))
		},
	}
}

func (app *AppHTTP) GetHTTPComponentConfig() Config.HTTP {
	return Config.HTTP{
		Server: TcpServer.New(8080, "", ""),
	}
}
