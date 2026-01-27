package main

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/dayvillefire/pocsag-monitor/config"
	"github.com/dayvillefire/pocsag-router/obj"
	"github.com/gin-gonic/gin"
)

func InitApi(m *gin.Engine) {
	a := Api{}
	g := m.Group("/api")

	g.GET("/debug", a.Debug)
	g.GET("/config", a.GetConfig)
	g.GET("/test/page/:capcode/:msg", a.TestPage)
	g.GET("/version", a.Version)
	/*
		g.GET("/calls/active", a.ActiveCalls)
		g.GET("/call/detail/:fdid/:id", a.CallDetail)
		g.GET("/version", a.Version)
		g.GET("/status/display", a.StatusDisplay)
	*/
}

type Api struct {
}

func (a Api) Debug(c *gin.Context) {
	o := map[string]any{}

	o["remote-ip"] = c.Request.RemoteAddr
	o["environment"] = os.Environ()
	o["user-id"] = os.Geteuid()
	o["pid"] = os.Getpid()
	o["architecture"] = runtime.GOARCH
	o["operating-system"] = runtime.GOOS
	//o["max-processes"] = runtime.GOMAXPROCS
	o["num-cpus"] = runtime.NumCPU()
	o["running-goroutines"] = runtime.NumGoroutine()

	c.JSON(http.StatusOK, o)
}

func (a Api) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.GetConfig())
}

func (a Api) TestPage(c *gin.Context) {
	capcode := c.Param("capcode")
	text := c.Param("msg")

	msg := obj.AlphaMessage{
		Timestamp: time.Now(),
		CapCode:   capcode,
		Message:   text,
		Valid:     true,
	}
	router.Publish(cfg.Router.Topic, msg)
	c.JSON(http.StatusOK, true)
}

func (a Api) Version(c *gin.Context) {
	c.JSON(http.StatusOK, Version)
}
