package handler

import "github.com/julienschmidt/httprouter"

type Hand interface {
	Register(router *httprouter.Router)
}
