module SystemgeSampleOauth2

go 1.23

//replace github.com/neutralusername/Systemge => ../Systemge

require (
	github.com/neutralusername/Systemge v0.0.0-20240825065049-29ea3104107e
	golang.org/x/oauth2 v0.21.0
)

require github.com/gorilla/websocket v1.5.3
