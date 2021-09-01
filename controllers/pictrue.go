package controllers

import (
	"github.com/kataras/iris/v12"
	"picture/protocal/resp"
	"picture/selenium"
)

func GetImages(context iris.Context) {
	count := context.URLParamIntDefault("count", 10)
	if count > 10 {
		context.JSON(resp.Fail(40000, "数量不能超过10"))
		return
	}
	keyword := context.URLParamDefault("keyword", "图片")
	imageSize := context.URLParam("imageSize")
	color := context.URLParam("color")
	imageType := context.URLParam("type")
	layout := context.URLParam("layout")

	result := selenium.GetImages(keyword, imageSize, color, imageType, layout, count)
	context.JSON(resp.Success(result))
}
