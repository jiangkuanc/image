package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	selenium2 "github.com/tebeka/selenium"
	"picture/controllers"
	"picture/selenium"
)

func main() {
	server := selenium.NewServer()
	defer func(server *selenium2.Service) {
		err := server.Stop()
		if err != nil {
		}
	}(server)

	app := iris.New()
	app.Get("/bing/images", controllers.GetImages)
	app.Configure(iris.WithConfiguration(iris.TOML("/data/server/image/iris.toml")))
	err := app.Run(iris.Addr(":9800"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
	if err != nil {
		return
	}
	fmt.Println("server stop .......................")
	server.Stop()
}
