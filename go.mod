module SystemgeSampleOauth2

go 1.23

//replace github.com/neutralusername/Systemge => ../Systemge

require (
	github.com/neutralusername/Systemge v0.0.0-20240908122750-9f1f787e8d01
	golang.org/x/oauth2 v0.21.0
)

require github.com/gorilla/websocket v1.5.3
