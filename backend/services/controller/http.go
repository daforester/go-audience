package controller

import (
	"github.com/gin-gonic/gin"
)

type HTTP struct {
	BaseController
	c       *gin.Context
	context *gin.Context
}

func (C *HTTP) GetContext() interface{} {
	return C.context
}

func (C *HTTP) SetContext(context interface{}) Controller {
	C.c = context.(*gin.Context)
	C.context = context.(*gin.Context)
	return C
}

func (C *HTTP) Run() error {
	return nil
}
