module github.com/dayvillefire/pocsag-monitor/controllers

go 1.23

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/config => ../config
	github.com/jbuchbinder/shims/factory => ../../../jbuchbinder/shims/factory
)

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/dayvillefire/pocsag-monitor/config v0.0.0-20241231155749-74823eee0d9a
	github.com/jbuchbinder/shims/factory v0.0.0-20240506232043-4fac4ec97ccb
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
