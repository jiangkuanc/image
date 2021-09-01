package selenium

import "time"

type SeleniumConfig struct {
	// 服务器监听端口
	Port int

	// google driver 路径
	DriverPath string

	// 打开标签页数量
	BookmarkPage int

	// 超时时间 (毫秒)
	Timeout time.Duration

	// bing 图片地址
	ImageUrl string
}
