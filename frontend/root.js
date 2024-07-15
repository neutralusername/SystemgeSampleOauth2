export class root extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            username: "",
            WS_CONNECTION: new WebSocket("ws://localhost:8443/ws"),
            constructMessage: (topic, payload) => {
                return JSON.stringify({
                    topic: topic,
                    payload: payload,
                });
            },
            setStateRoot: (state) => {
                this.setState(state)
            }
        },
        (this.state.WS_CONNECTION.onmessage = (event) => {
            let message = JSON.parse(event.data);
            switch (message.topic) {
                case "authRequest": {
                    let params = new URL(document.location.toString()).searchParams;
                    let sessionId = params.get("sessionId");
                    if (!sessionId) {
                        let cookies = document.cookie.split("; ");
                        for (let i = 0; i < cookies.length; i++) {
                            let cookie = cookies[i].split("=");
                            if (cookie[0] == "sessionId") {
                                sessionId = cookie[1];
                                break;
                            }
                        }
                    }
                    this.state.WS_CONNECTION.send(this.state.constructMessage("authAttempt", sessionId));
                    break;
                }
                case "authSuccess": {
                    let params = new URL(document.location.toString()).searchParams;
                    let sessionId = params.get("sessionId");
                    if (sessionId) {
                        document.cookie = "sessionId=" + sessionId + "; path=/";
                    }
                    window.history.replaceState({}, document.title, "/");
                    this.state.setStateRoot({
                        username : message.payload,
                    })
                    break;
                }
                case "authFailure":
                    document.cookie = "sessionId=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
                    window.history.replaceState({}, document.title, "/");
                    this.state.setStateRoot({
                        username : "",
                    })
                    break;
                case "logoutSuccess":
                    document.cookie = "sessionId=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT";
                    this.state.setStateRoot({
                        username : "",
                    })
                    break;
                case "logoutFailure":
                    console.log("Logout failed");
                    break;
                default:
                    console.log("Unknown message topic: " + event.data);
                    break;
            }
        });
        this.state.WS_CONNECTION.onclose = () => {
            setTimeout(() => {
                if (this.state.WS_CONNECTION.readyState === WebSocket.CLOSED) {}
                window.location.reload();
            }, 2000);
        };
        this.state.WS_CONNECTION.onopen = () => {
            let myLoop = () => {
                this.state.WS_CONNECTION.send(this.state.constructMessage("heartbeat", ""));
                setTimeout(myLoop, 15 * 1000);
            };
            setTimeout(myLoop, 15 * 1000);
        };
    }

    render() {
        return React.createElement(
            "div", {
                id: "root",
                onContextMenu: (e) => {
                    e.preventDefault();
                },
                style: {
                    fontFamily: "sans-serif",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                    touchAction: "none",
                    userSelect: "none",
                },
            },
            this.state.username == "" ? "unauthorized" : ("Hello, " + this.state.username),
            this.state.username == ""? React.createElement("button", {
                onClick: () => {
                    window.location.href = "http://localhost:8081";
                },
                style: {
                    marginTop: "10px",
                    padding: "5px",
                    backgroundColor: "white",
                    border: "1px solid black",
                    borderRadius: "5px",
                    cursor: "pointer",
                },
            }, "authorize") : 
            React.createElement("button", {
                onClick: () => {
                    let cookies = document.cookie.split("; ");
                    let sessionId = "";
                    for (let i = 0; i < cookies.length; i++) {
                        let cookie = cookies[i].split("=");
                        if (cookie[0] == "sessionId") {
                            sessionId = cookie[1];
                            break;
                        }
                    }
                    this.state.WS_CONNECTION.send(this.state.constructMessage("logoutAttempt", sessionId));
                },
                style: {
                    marginTop: "10px",
                    padding: "5px",
                    backgroundColor: "white",
                    border: "1px solid black",
                    borderRadius: "5px",
                    cursor: "pointer",
                },
            }, "logout")
        );
    }
}