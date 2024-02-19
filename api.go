package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/dayvillefire/pocsag-monitor/config"
	"github.com/dayvillefire/pocsag-monitor/output"
	"github.com/gin-gonic/gin"
)

func InitApi(m *gin.Engine) {
	a := Api{}
	g := m.Group("/api")

	g.GET("/debug", a.Debug)
	g.GET("/config", a.GetConfig)
	g.GET("/config/reload", a.ConfigReload)
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
	o := map[string]interface{}{}

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

func (a Api) ConfigReload(c *gin.Context) {
	d, err := config.ReloadDynamicConfig(*dynamicConfigFile)
	if err != nil {
		log.Print(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cfg.Dynamic = &d
	log.Printf("ConfigReload: Reloading dynamic as : %#v", cfg.Dynamic)

	// Dynamic channel mapping init
	routerMutex.Lock()
	defer routerMutex.Unlock()
	router = Router{cfg.Dynamic.ChannelMappings}
	outputs = map[string]output.Output{}
	for k, v := range cfg.Dynamic.OutputChannels {
		outputs[k], err = output.InstantiateOutput(v.Plugin)
		if err != nil {
			log.Printf(k + "| ERR: " + err.Error())
			c.JSON(http.StatusInternalServerError, k+": "+err.Error())
			return
		}
		err = outputs[k].Init(v.Option)
		if err != nil {
			log.Printf(k + "| ERR: " + err.Error())
			c.JSON(http.StatusInternalServerError, k+": "+err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, true)
}
