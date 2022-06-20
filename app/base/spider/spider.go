package spider


type (
	Spider struct {
		Name 				string
		Description 		string

		// 采集规则树


	}

	// 采集规则树
	RuleTree struct {
		Root 	func()									// 根结点(执行入口)
		Trunk	map[string]*Rule						// 节点散列表(采集过程)
	}

	Rule struct {
		ItemField 	[]string
		ParseFunc 	func()
		AidFunc		func(interface{}) interface{}		// 辅助函数
	}
)
