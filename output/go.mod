module github.com/dayvillefire/pocsag-monitor/output

go 1.20

replace github.com/dayvillefire/pocsag-monitor/obj => ../obj

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20230412213233-faa2c43219a4
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
)
