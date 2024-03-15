module github.com/dayvillefire/pocsag-monitor/output

go 1.22

replace github.com/dayvillefire/pocsag-monitor/obj => ../obj

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20240107151918-75d6209d56ee
	github.com/eclipse/paho.mqtt.golang v1.4.3
)

require (
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)
