module github.com/dayvillefire/pocsag-monitor/output

go 1.22

replace github.com/dayvillefire/pocsag-monitor/obj => ../obj

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20240703200132-b283f52aa359
	github.com/eclipse/paho.mqtt.golang v1.4.3
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
)
