package controllers

import (
	"github.com/doubletrey/crawlab-core/entity"
	"github.com/doubletrey/crawlab-core/interfaces"
	"github.com/gin-gonic/gin"
)

type BinderInterface interface {
	Bind(c *gin.Context) (res interfaces.Model, err error)
	BindList(c *gin.Context) (res []interfaces.Model, err error)
	BindBatchRequestPayload(c *gin.Context) (payload entity.BatchRequestPayload, err error)
	BindBatchRequestPayloadWithStringData(c *gin.Context) (payload entity.BatchRequestPayloadWithStringData, res interfaces.Model, err error)
}
