package controllers

import (
	"github.com/doubletrey/crawlab-core/constants"
	"github.com/doubletrey/crawlab-core/entity"
	"github.com/gin-gonic/gin"
)

func GetDefaultPagination() (p *entity.Pagination) {
	return &entity.Pagination{
		Page: constants.PaginationDefaultPage,
		Size: constants.PaginationDefaultSize,
	}
}

func GetPagination(c *gin.Context) (p *entity.Pagination, err error) {
	var _p entity.Pagination
	if err := c.ShouldBindQuery(&_p); err != nil {
		return GetDefaultPagination(), err
	}
	return &_p, nil
}

func MustGetPagination(c *gin.Context) (p *entity.Pagination) {
	p, err := GetPagination(c)
	if err != nil || p == nil {
		return GetDefaultPagination()
	}
	return p
}
