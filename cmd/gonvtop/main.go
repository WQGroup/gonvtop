package main

import (
	"errors"
	"fmt"
	"github.com/WQGroup/gonvtop/pkg/cors"
	v1 "github.com/WQGroup/gonvtop/pkg/http_api/v1"
	"github.com/WQGroup/gonvtop/pkg/info_hub"
	"github.com/WQGroup/logger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func main() {

	// ----------------- Gin -----------------
	
	var srv *http.Server
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	engine := gin.Default()
	engine.ForwardedByClientIP = true
	// 默认所有都通过
	engine.Use(cors.Cors())
	// 后端实例
	infoHub := info_hub.NewInfoHub("")
	defer infoHub.Close()
	cbV1 := v1.NewControllerBase(infoHub)

	groupV1 := engine.Group("/infos/" + cbV1.GetVersion())
	{
		groupV1.GET("/gpu_driver", cbV1.GetGPUDriverInfos)

		groupV1.GET("/gpus", cbV1.GetGPUs)

		groupV1.GET("/host", cbV1.GetHostInfos)
	}

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", 19035),
		Handler: engine,
	}

	logger.Infoln("Try Start Http Server At Port", "19035")
	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) == false {
		logger.Panicln("Start Server Error:", err)
	}
}
