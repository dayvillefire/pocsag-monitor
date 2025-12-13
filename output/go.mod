module github.com/dayvillefire/pocsag-monitor/output

go 1.25

replace github.com/dayvillefire/pocsag-monitor/obj => ../obj

replace github.com/dayvillefire/groupme => ../../groupme

require (
	github.com/bwmarrin/discordgo v0.29.0
	github.com/dayvillefire/groupme v0.5.1
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20251208221133-403e471ee4eb
	github.com/eclipse/paho.mqtt.golang v1.5.1
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
)
