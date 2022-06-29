package config

import (
	"skynet-service/app/common"
	"strings"

	"skynet-service/app/logs/logs"
	"skynet-service/app/runtime/cache"
)

// 软件信息。
const (
	VERSION  string = "v0.0.0"                                      		// 软件版本号
	AUTHOR   string = "DingJing"                                       		// 软件作者
	NAME     string = "spider"                                 				// 软件名
	FullName string = NAME + " " + VERSION + "（by " + AUTHOR + "）" 		// 软件全称
	IconPng  string = ``
	PidFile	 string = "/tmp/" + NAME + "_" + AUTHOR + "_" + VERSION			// 进程id存放位置
)

// 默认配置。
const (
	WorkRoot      string = common.TAG + "_running"               // 运行时的目录名称
	CONFIG        string = WorkRoot + "/config.ini"       // 配置文件路径
	CacheDir      string = WorkRoot + "/cache"            // 缓存文件目录
	LOG           string = WorkRoot + "/logs/spiders.log" // 日志文件路径
	LogAsync      bool   = true                           // 是否异步输出日志
	PhantomjsTemp string = CacheDir                       // Surfer-Phantom下载器：js文件临时目录
	HistoryTag    string = "history"                      // 历史记录的标识符
	HistoryDir    string = WorkRoot + "/" + HistoryTag    // excel或csv输出方式下，历史记录目录
	SpiderExt     string = ".spider.html"                 // 动态规则扩展名
)

// 来自配置文件的配置项。
var (
	CRAWLS_CAP int = setting.DefaultInt("crawlcap", crawlcap) // 蜘蛛池最大容量
	// DATA_CHAN_CAP            int    = setting.DefaultInt("datachancap", datachancap)                         // 收集器容量
	PHANTOMJS                string = setting.String("phantomjs")                                          // Surfer-Phantom下载器：phantomjs程序路径
	PROXY                    string = setting.String("proxylib")                                           // 代理IP文件路径
	SPIDER_DIR               string = setting.String("spiderdir")                                          // 动态规则目录
	FILE_DIR                 string = setting.String("fileoutdir")                                         // 文件（图片、HTML等）结果的输出目录
	TEXT_DIR                 string = setting.String("textoutdir")                                         // excel或csv输出方式下，文本结果的输出目录
	DB_NAME                  string = setting.String("dbname")                                             // 数据库名称
	MGO_CONN_STR             string = setting.String("mgo::connstring")                                    // mongodb连接字符串
	MGO_CONN_CAP             int    = setting.DefaultInt("mgo::conncap", mgoconncap)                       // mongodb连接池容量
	MGO_CONN_GC_SECOND       int64  = setting.DefaultInt64("mgo::conngcsecond", mgoconngcsecond)           // mongodb连接池GC时间，单位秒
	MYSQL_CONN_STR           string = setting.String("mysql::connstring")                                  // mysql连接字符串
	MYSQL_CONN_CAP           int    = setting.DefaultInt("mysql::conncap", mysqlconncap)                   // mysql连接池容量
	MYSQL_MAX_ALLOWED_PACKET int    = setting.DefaultInt("mysql::maxallowedpacket", mysqlmaxallowedpacket) // mysql通信缓冲区的最大长度
	BeanstalkdHost           string = setting.DefaultString("beanstalkd::host", beanstalkHost)             // Beanstalkd指定主机地址
	BeanstalkdTube           string = setting.DefaultString("beanstalkd::tube", beanstalkTube)             // Beanstalkd指定主机地址

	KAFKA_BORKERS string = setting.DefaultString("kafka::brokers", kafkabrokers) 							// kafka brokers

	LOG_CAP            int64 = setting.DefaultInt64("log::cap", logcap)          // 日志缓存的容量
	LOG_LEVEL          int   = logLevel(setting.String("log::level"))            // 全局日志打印级别（亦是日志文件输出级别）
	LOG_CONSOLE_LEVEL  int   = logLevel(setting.String("log::consolelevel"))     // 日志在控制台的显示级别
	LOG_FEEDBACK_LEVEL int   = logLevel(setting.String("log::feedbacklevel"))    // 客户端反馈至服务端的日志级别
	LOG_LINEINFO       bool  = setting.DefaultBool("log::lineinfo", loglineinfo) // 日志是否打印行信息                                  // 客户端反馈至服务端的日志级别
	LOG_SAVE           bool  = setting.DefaultBool("log::save", logsave)         // 是否保存所有日志到本地文件
)

func init() {
	// 主要运行时参数的初始化
	cache.Task = &cache.AppConf{
		Mode:           setting.DefaultInt("run::mode", mode),                 // 节点角色
		Port:           setting.DefaultInt("run::port", port),                 // 主节点端口
		Master:         setting.String("run::master"),                         // 服务器(主节点)地址，不含端口
		ThreadNum:      setting.DefaultInt("run::thread", thread),             // 全局最大并发量
		Pausetime:      setting.DefaultInt64("run::pause", pause),             // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
		OutType:        setting.String("run::outtype"),                        // 输出方式
		DockerCap:      setting.DefaultInt("run::dockercap", dockercap),       // 分段转储容器容量
		Limit:          setting.DefaultInt64("run::limit", limit),             // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
		ProxyMinute:    setting.DefaultInt64("run::proxyminute", proxyminute), // 代理IP更换的间隔分钟数
		SuccessInherit: setting.DefaultBool("run::success", success),          // 继承历史成功记录
		FailureInherit: setting.DefaultBool("run::failure", failure),          // 继承历史失败记录
	}
}

func logLevel(l string) int {
	switch strings.ToLower(l) {
	case "app":
		return logs.LevelApp
	case "emergency":
		return logs.LevelEmergency
	case "alert":
		return logs.LevelAlert
	case "critical":
		return logs.LevelCritical
	case "error":
		return logs.LevelError
	case "warning":
		return logs.LevelWarning
	case "notice":
		return logs.LevelNotice
	case "informational":
		return logs.LevelInformational
	case "info":
		return logs.LevelInformational
	case "debug":
		return logs.LevelDebug
	}
	return -10
}
