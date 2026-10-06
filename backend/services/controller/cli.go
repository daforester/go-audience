package controller

import (
	"github.com/daforester/getopt-golang/getopt"
)

type CLI struct {
	BaseController
	c       *getopt.GetOpt
	context *getopt.GetOpt
}

func (C *CLI) GetContext() interface{} {
	return C.context
}

func (C *CLI) SetContext(context interface{}) Controller {
	C.c = context.(*getopt.GetOpt)
	C.context = context.(*getopt.GetOpt)
	return C
}

func (C *CLI) Run() error {
	return nil
}
