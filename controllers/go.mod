module github.com/dayvillefire/pocsag-monitor/controllers

go 1.22

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/config => ../config
	github.com/jbuchbinder/shims/factory => ../../../jbuchbinder/shims/factory
)

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/dayvillefire/pocsag-monitor/config v0.0.0-20240718135929-e5e4f1babde7
	github.com/jbuchbinder/shims/factory v0.0.0-20240506232043-4fac4ec97ccb
)

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
