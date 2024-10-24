module SystemgeSampleOauth2

go 1.23

replace github.com/neutralusername/Systemge => ../Systemge

require (
	github.com/neutralusername/Systemge v0.0.0-20240920150811-762a862539cc
	golang.org/x/oauth2 v0.21.0
)

require github.com/gorilla/websocket v1.5.3

require github.com/neutralusername/systemge v0.0.0-20241024165516-386f4dc5f91c // indirect
