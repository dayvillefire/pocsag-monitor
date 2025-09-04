module github.com/dayvillefire/pocsag-monitor/controllers

go 1.23.0

toolchain go1.24.2

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/config => ../config
	github.com/jbuchbinder/shims/factory => ../../../jbuchbinder/shims/factory
)

require (
	github.com/bwmarrin/discordgo v0.29.0
	github.com/dayvillefire/pocsag-monitor/config v0.0.0-20250828124955-f581123ea150
	github.com/jbuchbinder/shims/factory v0.0.0-20250818154854-22c0ac83b788
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
