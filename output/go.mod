module github.com/dayvillefire/pocsag-monitor/output

go 1.23.0

toolchain go1.24.2

replace github.com/dayvillefire/pocsag-monitor/obj => ../obj

require (
	github.com/bwmarrin/discordgo v0.29.0
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20250411125416-eecc4a8a5f6c
	github.com/eclipse/paho.mqtt.golang v1.5.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)
