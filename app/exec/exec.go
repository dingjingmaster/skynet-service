package exec

import (
	"runtime"
	"skynet-service/app"
	"skynet-service/app/common/gc"
	"skynet-service/app/spider"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU()) 				// 开启最大核心数运行
	gc.ManualGC()										// 开启手动GC
}

func Run () {
	// 启动网页显示服务

	// 启动爬虫任务
	RunSpider()

	// 启动所有 spider
	//for true {
	//	;
	//}

	// 开始运行
	//ctrl := make(chan os.Signal, 1)
	//signal.Notify(ctrl, os.Interrupt, os.Kill)
	////go web.Run()
	//<-ctrl
}

// FIXME:// 此处需要根据命令行进行配置
// 启动 spider
func RunSpider() {
	app.LogicApp.Init()

	GetAllSpider()

	app.LogicApp.Run()
}

// 扫描爬虫
func GetAllSpider () {
	var spiders []*spider.Spider
	for _, sp := range app.LogicApp.GetSpiderLib() {
		spiders = append(spiders, sp.Copy())
	}
	app.LogicApp.SpiderPrepare(spiders)
}