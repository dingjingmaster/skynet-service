package pipeline

import (
	"skynet-service/app/common/kafka"
	"skynet-service/app/common/mgo"
	"skynet-service/app/common/mysql"
	"skynet-service/app/pipeline/collector"
	"skynet-service/app/runtime/cache"
	"sort"
)

// 初始化输出方式列表collector.DataOutputLib
func init() {
	for out, _ := range collector.DataOutput {
		collector.DataOutputLib = append(collector.DataOutputLib, out)
	}
	sort.Strings(collector.DataOutputLib)
}

// 刷新输出方式的状态
func RefreshOutput() {
	switch cache.Task.OutType {
	case "mgo":
		mgo.Refresh()
	case "mysql":
		mysql.Refresh()
	case "kafka":
		kafka.Refresh()
	}
}
