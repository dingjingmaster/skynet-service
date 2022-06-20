package spider

import (
	"github.com/gocolly/colly"
	"skynet-service/app/common/status"
	"skynet-service/app/spiders"
	"sync"
)

type (
	Spider struct {
		Name 				string
		Description 		string
		RuleTree			*RuleTree

		RequestType 		int							// 请求类型: POST/GET

		//
		status				int							// 运行状态
		lock 				sync.RWMutex
		once 				sync.Once
	}

	// 采集规则树
	RuleTree struct {
		RootURI			string
		Trunk 			map[string]*Rule				// 节点散列表(采集过程)
	}

	Rule struct {
		ItemFields 		[]string
		RequestFunc 	func(r *colly.Request)
		ResponseFunc 	func(r *colly.Response)

		HTMLParser 		map[string]*colly.HTMLElement
		XMLParser 		map[string]*colly.XMLElement

		AidFunc			func(interface{}) interface{}	// 辅助函数
	}
)


func (self Spider) Register() *Spider {
	self.status = status.SPIDER_STOPPED

	return spiders.Species.Add(&self)
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

		ghost.RuleTree.Trunk[k].RequestFunc = v.RequestFunc
		ghost.RuleTree.Trunk[k].ResponseFunc = v.ResponseFunc
		ghost.RuleTree.Trunk[k].HTMLParser = v.HTMLParser
		ghost.RuleTree.Trunk[k].XMLParser = v.XMLParser
		ghost.RuleTree.Trunk[k].AidFunc = v.AidFunc
	}

	ghost.Description = self.Description
	ghost.status = self.status

	return ghost
}

func (self *Spider) defaultRootRequest (r map[string]string) {
}

func New () {

}

// 开始执行
func (self *Spider) Start () {
	defer func() {
		self.lock.Lock()
		self.status = status.SPIDER_RUNNING
		self.lock.Unlock()
	}()

	//ctx = colly.NewCollector()
}

// 退出任务前收尾工作
func (self *Spider) Defer() {

}