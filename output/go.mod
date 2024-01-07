module github.com/dayvillefire/pocsag-monitor/output

go 1.21

replace github.com/dayvillefire/pocsag-monitor/obj => ../obj

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20230902204846-5d7fdba5601a
)

require (
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
