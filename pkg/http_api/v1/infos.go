package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (cb ControllerBase) GetGPUDriverInfos(c *gin.Context) {

	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "GetGPUDriverInfos", err)
	}()

}

func (cb ControllerBase) GetGPUs(c *gin.Context) {
	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "GetGPUs", err)
	}()

}

func (cb ControllerBase) GetHostInfos(c *gin.Context) {

	var err error
	defer func() {
		// 统一的异常处理
		cb.ErrorProcess(c, "GetHostInfos", err)
	}()

	c.JSON(http.StatusOK, cb.infoHub.GetCacheInfo())
}
