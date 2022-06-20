package spiders

import (
	"fmt"
	"skynet-service/app/base/spider"
	"spider/common/pinyin"
)

// spider 类型/列表
type SpiderSpecies struct {
	list 		[]*spider.Spider
	hash 		map[string]*spider.Spider
	sorted		bool
}

var Species = &SpiderSpecies {
	list: []*spider.Spider{},
	hash: map[string]*spider.Spider{},
}

func (self* SpiderSpecies) Add (sp* spider.Spider) *spider.Spider {
	name := sp.Name
	for i := 2; true; i++ {
		if _, ok := self.hash[name]; !ok {
			sp.Name = name
			self.hash[sp.Name] = sp
			break
		}
		name = fmt.Sprint("%s(%d)", sp.Name, i)
	}
	sp.Name = name
	self.list = append(self.list, sp)

	return sp
}


func (self *SpiderSpecies) Get() []*spider.Spider {
	if !self.sorted {
		l := len(self.list)
		initials := make([]string, l)
		newList := map[string]*spider.Spider{}
		for i := 0; i < l; i++ {
			initials[i] = self.list[i].GetName()
			newList[initials[i]] = self.list[i]
		}
		pinyin.SortInitials(initials)
		for i := 0; i < l; i++ {
			self.list[i] = newList[initials[i]]
		}
		self.sorted = true
	}

	return self.list
}

func (self *SpiderSpecies) GetByName (name string) *spider.Spider {
	return self.hash[name]
}
}