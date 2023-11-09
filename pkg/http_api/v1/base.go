package v1

import (
	"github.com/WQGroup/gonvtop/pkg/info_hub"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ControllerBase struct {
	infoHub *info_hub.InfoHub
}

func NewControllerBase(infoHub *info_hub.InfoHub) *ControllerBase {
	go func() {
		infoHub.Monitor(1 * time.Second)
	}()
	return &ControllerBase{infoHub: infoHub}
}

func (cb ControllerBase) GetVersion() string {
	return "v1"
}

func (cb ControllerBase) ErrorProcess(c *gin.Context, funcName string, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, ReplyCommon{Message: funcName + ":" + err.Error()})
	}
}

type ReplyCommon struct {
	Message string `json:"message,omitempty"`
}
