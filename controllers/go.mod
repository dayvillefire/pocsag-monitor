module github.com/dayvillefire/pocsag-monitor/controllers

go 1.22

replace github.com/jbuchbinder/shims/factory => ../../../jbuchbinder/shims/factory

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/jbuchbinder/shims/factory v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
