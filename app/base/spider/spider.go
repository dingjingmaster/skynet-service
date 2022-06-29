package spider

import (
	"skynet-service/app/common/status"
	"spider/app/downloader/request"
	"spider/app/scheduler"
	"sync"
)

type (
	Spider struct {
		Name 				string
		Description 		string
		RuleTree			*RuleTree

		RequestType 		int												// 请求类型: POST/GET

		//
		status				int												// 运行状态
		reqMatrix 			*scheduler.Matrix 								// 请求矩阵
		lock 				sync.RWMutex
		once 				sync.Once
	}

	// 采集规则树
	RuleTree struct {
		Root			func(c *Context)									// 根节点入口
		Trunk 			map[string]*Rule									// 节点散列表(采集过程)
	}

	Rule struct {
		ItemFields 		[]string
		ParseFunc  		func(*Context)                  					// 内容解析函数
		AidFunc			func(*Context, map[string]interface{}) interface{}	// 辅助函数
	}
)


func (self Spider) Register() *Spider {
	self.status = status.SPIDER_STOPPED

	return Species.Add(&self)
}

// 指定规则的获取结果的字段名列表
func (self *Spider) GetItemFields(rule *Rule) []string {
	return rule.ItemFields
}

// 返回结果字段名的值, 不存在时返回空字符串
func (self *Spider) GetItemField(rule *Rule, index int) (field string) {
	if index > len(rule.ItemFields)-1 || index < 0 {
		return ""
	}
	return rule.ItemFields[index]
}

// 返回结果字段名的其索引, 不存在时索引为-1
func (self *Spider) GetItemFieldIndex(rule *Rule, field string) (index int) {
	for idx, v := range rule.ItemFields {
		if v == field {
			return idx
		}
	}
	return -1
}

// 为指定Rule动态追加结果字段名，并返回索引位置, 已存在时返回原来索引位置
func (self *Spider) UpsertItemField(rule *Rule, field string) (index int) {
	for i, v := range rule.ItemFields {
		if v == field {
			return i
		}
	}
	rule.ItemFields = append(rule.ItemFields, field)
	return len(rule.ItemFields) - 1
}

// 获取蜘蛛名称
func (self *Spider) GetName() string {
	return self.Name
}

// 安全返回指定规则
func (self *Spider) GetRule(ruleName string) (*Rule, bool) {
	rule, found := self.RuleTree.Trunk[ruleName]
	return rule, found
}

// 返回指定规则
func (self *Spider) MustGetRule(ruleName string) *Rule {
	return self.RuleTree.Trunk[ruleName]
}

// 返回规则树
func (self *Spider) GetRules() map[string]*Rule {
	return self.RuleTree.Trunk
}

// 获取蜘蛛描述
func (self *Spider) GetDescription() string {
	return self.Description
}

// 返回一个自身复制品
func (self *Spider) Copy() *Spider {
	ghost := &Spider{}
	ghost.Name = self.Name

	ghost.RuleTree = &RuleTree{
		Trunk: make(map[string]*Rule, len(self.RuleTree.Trunk)),
	}
	for k, v := range self.RuleTree.Trunk {
		ghost.RuleTree.Trunk[k] = new(Rule)

		ghost.RuleTree.Trunk[k].ItemFields = make([]string, len(v.ItemFields))
		copy(ghost.RuleTree.Trunk[k].ItemFields, v.ItemFields)

		ghost.RuleTree.Trunk[k].ParseFunc = v.ParseFunc
		ghost.RuleTree.Trunk[k].AidFunc = v.AidFunc
	}

	ghost.Description = self.Description
	ghost.status = self.status

	return ghost
}

func (self *Spider) defaultRootRequest (r map[string]string) {
}

func (self *Spider) RequestPush(req *request.Request) {
	self.reqMatrix.Push(req)
}

// 开始执行
func (self *Spider) Start () {
	defer func() {
		self.lock.Lock()
		self.status = status.SPIDER_RUNNING
		self.lock.Unlock()
	}()

	self.RuleTree.Root(GetContext(self, nil))
}

// 主动崩溃爬虫运行协程
func (self *Spider) Stop () {
	self.lock.Lock()

	defer self.lock.Unlock()

	if self.status == status.SPIDER_STOPPED {
		return
	}
	self.status = status.SPIDER_STOPPED
}

// 退出任务前收尾工作
func (self *Spider) Defer() {

}
