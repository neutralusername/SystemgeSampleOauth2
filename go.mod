module SystemgeSampleOauth2

go 1.23

//replace github.com/neutralusername/Systemge => ../Systemge

require (
	github.com/neutralusername/Systemge v0.0.0-20240907083548-aec109fbc08c
	golang.org/x/oauth2 v0.21.0
)

require github.com/gorilla/websocket v1.5.3
