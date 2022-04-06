package controllers

import (
	"github.com/doubletrey/crawlab-core/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetVersion(c *gin.Context) {
	HandleSuccessWithData(c, config.GetVersion())
}

func getVersionActions() []Action {
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: GetVersion,
		},
	}
}

var VersionController ActionController
