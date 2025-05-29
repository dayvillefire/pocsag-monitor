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
	github.com/dayvillefire/pocsag-monitor/config v0.0.0-20250411125416-eecc4a8a5f6c
	github.com/jbuchbinder/shims/factory v0.0.0-20250315180801-ea13cafaf717
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
