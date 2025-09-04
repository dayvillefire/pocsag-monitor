module github.com/dayvillefire/pocsag-monitor/output

go 1.23.0

toolchain go1.24.3

replace github.com/dayvillefire/pocsag-monitor/obj => ../obj

replace github.com/dayvillefire/groupme => ../../groupme

require (
	github.com/bwmarrin/discordgo v0.29.0
	github.com/dayvillefire/groupme v0.5.1
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20250828124955-f581123ea150
	github.com/eclipse/paho.mqtt.golang v1.5.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
