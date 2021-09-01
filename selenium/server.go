package selenium

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"os"
	"strings"
)

var (
	config    = new(SeleniumConfig)
	service   = selenium.Service{}
	// WebDriver = make([]selenium.WebDriver, 0)
)

func init() {
	if len(os.Args) == 1 {
		panic("未指定环境 可选值 dev、test、prod")
	} else if len(os.Args) > 2 {
		panic("参数不正确 可选值 dev、test、prod")
	}
	configType := os.Args[1]
	if _, err := toml.DecodeFile("/data/server/image/selenium-config-"+configType+".toml", &config); err != nil {
		log.Fatal(err)
	}
}

func NewServer() *selenium.Service {
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	service, err := selenium.NewChromeDriverService(config.DriverPath, config.Port, opts...)
	if err != nil {
		return nil
	}
	// caps := selenium.Capabilities{"browserName": "chrome"}
	// for i := 0; i < config.BookmarkPage; i++ {
	// 	var windowOne, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", config.Port))
	// 	if err != nil {
	// 		print("get windowOne error")
	// 	}
	// 	WebDriver = append(WebDriver, windowOne)
	// }
	return service
}

func GetImages(keyword, imageSize, color, imageType, layout string, count int) []string {
	var condition []string
	if keyword != "" {
		condition = append(condition, fmt.Sprintf("q=%s", keyword))
	}
	var filterUI string
	if imageSize != "" || color != "" || imageType != "" || layout != "" {
		filterUI = "qft="
		if imageSize != "" {
			filterUI += fmt.Sprintf("+filterui:imagesize-%s", strings.ToLower(imageSize))
		}
		if color != "" {
			filterUI += fmt.Sprintf("+filterui:color2-FGcls_%s", strings.ToTitle(color))
		}
		if imageType != "" {
			filterUI += fmt.Sprintf("+filterui:photo-%s", strings.ToLower(imageType))
		}
		if layout != "" {
			filterUI += fmt.Sprintf("+filterui:aspect-%s", strings.ToLower(layout))
		}
	}
	if filterUI != "" {
		condition = append(condition, filterUI)
	}
	var url = config.ImageUrl
	if len(condition) > 0 {
		url += "&" + strings.Join(condition, "&")
	}

	// TODO 获取空闲的连接，超过最大等待时间返回 服务器繁忙，请稍后重试
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			// 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--headless",
			"--no-sandbox",
			// 模拟user-agent，防反爬
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	windowOne, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", config.Port))
	defer  windowOne.Quit()

	_ = windowOne.Get(url)

	fmt.Printf("request url:%s", url)

	// Get a reference to the text box containing code.
	componentImages, err := windowOne.FindElements(selenium.ByXPATH, "//*[@class=\"dgControl hover\"]/ul")
	if err != nil {
		panic(err)
	}
	result := make([]string, 0)
	if len(componentImages) > 0 {
		for i := 0; i < len(componentImages); i++ {
			img, err := componentImages[i].FindElements(selenium.ByTagName, "img")
			if err != nil {
				continue
			}
			if len(img) > 0 {
				for j := 0; j < len(img); j++ {
					imgUrl, err := img[j].GetAttribute("src")
					if err != nil {
						continue
					}
					if len(result) == count {
						break
					} else {
						result = append(result, imgUrl)
					}
				}
			}
			if len(result) == count {
				break
			}
		}
	}
	return result
}
