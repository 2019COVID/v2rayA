package controller

import (
	"V2RayA/global"
	"V2RayA/persistence/configure"
	"V2RayA/service"
	"V2RayA/tools"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func GetPortWhiteList(ctx *gin.Context) {
	tools.ResponseSuccess(ctx, configure.GetPortWhiteListNotNil())
}

func PutPortWhiteList(ctx *gin.Context) {
	var data configure.PortWhiteList
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		tools.ResponseError(ctx, errors.New("bad request"))
		return
	}
	if !data.Valid() {
		tools.ResponseError(ctx, errors.New("invalid format of port"))
		return
	}
	err = configure.SetPortWhiteList(&data)
	if err != nil {
		tools.ResponseError(ctx, err)
		return
	}
	tools.ResponseSuccess(ctx, nil)
}

func PostPortWhiteList(ctx *gin.Context) {
	var data struct {
		RequestPort string `json:"requestPort"`
	}
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		tools.ResponseError(ctx, errors.New("bad request"))
		return
	}
	p := service.GetPortsDefault()
	listenPort := strings.Split(global.GetEnvironmentConfig().Address, ":")[1]
	wpl := configure.PortWhiteList{
		TCP: []string{"1:1023", data.RequestPort, listenPort, strconv.Itoa(p.Http), strconv.Itoa(p.HttpWithPac), strconv.Itoa(p.Socks5)},
		UDP: []string{"1:1023", strconv.Itoa(p.Socks5)},
	}
	err = configure.SetPortWhiteList(&wpl)
	if err != nil {
		tools.ResponseError(ctx, err)
		return
	}
	tools.ResponseSuccess(ctx, nil)
}
