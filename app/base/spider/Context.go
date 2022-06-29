package spider

import (
	"bytes"
	"golang.org/x/net@v0.0.0-20220617184016-355a448f1bc9/html/charset"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"spider/app/downloader/request"
	"spider/app/pipeline/collector/data"
	"spider/common/goquery"
	"spider/common/util"
	"spider/logs"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type Context struct {
	spider   *Spider           			// 规则
	Request  *request.Request  			// 原始请求
	Response *http.Response    			// 响应流，其中URL拷贝自*request.Request
	text     []byte            			// 下载内容Body的字节流格式
	dom      *goquery.Document 			// 下载内容Body为html时，可转换为Dom的对象
	items    []data.DataCell   			// 存放以文本形式输出的结果数据
	files    []data.FileCell   			// 存放欲直接输出的文件("Name": string; "Body": io.ReadCloser)
	err      error             			// 错误标记
	sync.Mutex
}

var (
	contextPool = &sync.Pool {
		New: func() interface{} {
			return &Context{
				items: []data.DataCell{},
				files: []data.FileCell{},
			}
		},
	}
)

func GetContext(sp *Spider, req *request.Request) *Context {
	ctx := contextPool.Get().(*Context)
	ctx.spider = sp
	ctx.Request = req
	return ctx
}

func PutContext(ctx *Context) {
	if ctx.Response != nil {
		ctx.Response.Body.Close() // too many open files bug remove
		ctx.Response = nil
	}
	ctx.items = ctx.items[:0]
	ctx.files = ctx.files[:0]
	ctx.spider = nil
	ctx.Request = nil
	ctx.text = nil
	ctx.dom = nil
	ctx.err = nil
	contextPool.Put(ctx)
}

func (self *Context) SetResponse(resp *http.Response) *Context {
	self.Response = resp
	return self
}

// 标记下载错误。
func (self *Context) SetError(err error) {
	self.err = err
}

func (self *Context) AddQueue(req *request.Request) *Context {

	err := req.
		SetSpiderName(self.spider.GetName()).
		SetEnableCookie(true).
		Prepare()

	if err != nil {
		logs.Log.Error(err.Error())
		return self
	}

	// 自动设置Referer
	if req.GetReferer() == "" && self.Response != nil {
		req.SetReferer(self.GetUrl())
	}

	self.spider.RequestPush(req)
	return self
}

func (self *Context) SetUrl(url string) *Context {
	self.Request.Url = url
	return self
}

func (self *Context) SetReferer(referer string) *Context {
	self.Request.Header.Set("Referer", referer)
	return self
}

func (self *Context) Aid(aid map[string]interface{}, ruleName ...string) interface{} {

	_, rule, found := self.getRule(ruleName...)
	if !found {
		if len(ruleName) > 0 {
			logs.Log.Error("调用蜘蛛 %s 不存在的规则: %s", self.spider.GetName(), ruleName[0])
		} else {
			logs.Log.Error("调用蜘蛛 %s 的Aid()时未指定的规则名", self.spider.GetName())
		}
		return nil
	}
	if rule.AidFunc == nil {
		logs.Log.Error("蜘蛛 %s 的规则 %s 未定义AidFunc", self.spider.GetName(), ruleName[0])
		return nil
	}
	return rule.AidFunc(self, aid)
}

// 解析响应流。
// 用ruleName指定匹配的ParseFunc字段，为空时默认调用Root()。
func (self *Context) Parse(ruleName ...string) *Context {

	_ruleName, rule, found := self.getRule(ruleName...)
	if self.Response != nil {
		self.Request.SetRuleName(_ruleName)
	}
	if !found {
		self.spider.RuleTree.Root(self)
		return self
	}
	if rule.ParseFunc == nil {
		logs.Log.Error("蜘蛛 %s 的规则 %s 未定义ParseFunc", self.spider.GetName(), ruleName[0])
		return self
	}
	rule.ParseFunc(self)
	return self
}

// 重置下载的文本内容，
func (self *Context) ResetText(body string) *Context {
	x := (*[2]uintptr)(unsafe.Pointer(&body))
	h := [3]uintptr{x[0], x[1], x[1]}
	self.text = *(*[]byte)(unsafe.Pointer(&h))
	self.dom = nil
	return self
}

// 获取下载错误。
func (self *Context) GetError() error {
	return self.err
}

// 获取日志接口实例。
func (*Context) Log() logs.Logs {
	return logs.Log
}

// 获取蜘蛛名称。
func (self *Context) GetSpider() *Spider {
	return self.spider
}

// 获取响应流。
func (self *Context) GetResponse() *http.Response {
	return self.Response
}

// 获取响应状态码。
func (self *Context) GetStatusCode() int {
	return self.Response.StatusCode
}

// 获取原始请求。
func (self *Context) GetRequest() *request.Request {
	return self.Request
}

// 获得一个原始请求的副本。
func (self *Context) CopyRequest() *request.Request {
	return self.Request.Copy()
}

// 获取结果字段名列表。
func (self *Context) GetItemFields(ruleName ...string) []string {
	_, rule, found := self.getRule(ruleName...)
	if !found {
		logs.Log.Error("蜘蛛 %s 调用GetItemFields()时，指定的规则名不存在！", self.spider.GetName())
		return nil
	}
	return self.spider.GetItemFields(rule)
}

// 由索引下标获取结果字段名，不存在时获取空字符串，
// 若ruleName为空，默认为当前规则。
func (self *Context) GetItemField(index int, ruleName ...string) (field string) {
	_, rule, found := self.getRule(ruleName...)
	if !found {
		logs.Log.Error("蜘蛛 %s 调用GetItemField()时，指定的规则名不存在！", self.spider.GetName())
		return
	}
	return self.spider.GetItemField(rule, index)
}

// 由结果字段名获取索引下标，不存在时索引为-1，
// 若ruleName为空，默认为当前规则。
func (self *Context) GetItemFieldIndex(field string, ruleName ...string) (index int) {
	_, rule, found := self.getRule(ruleName...)
	if !found {
		logs.Log.Error("蜘蛛 %s 调用GetItemField()时，指定的规则名不存在！", self.spider.GetName())
		return
	}
	return self.spider.GetItemFieldIndex(rule, field)
}

func (self *Context) PullItems() (ds []data.DataCell) {
	self.Lock()
	ds = self.items
	self.items = []data.DataCell{}
	self.Unlock()
	return
}

func (self *Context) PullFiles() (fs []data.FileCell) {
	self.Lock()
	fs = self.files
	self.files = []data.FileCell{}
	self.Unlock()
	return
}

// 获取蜘蛛名。
func (self *Context) GetName() string {
	return self.spider.GetName()
}

// 获取规则树。
func (self *Context) GetRules() map[string]*Rule {
	return self.spider.GetRules()
}

// 获取指定规则。
func (self *Context) GetRule(ruleName string) (*Rule, bool) {
	return self.spider.GetRule(ruleName)
}

// 获取当前规则名。
func (self *Context) GetRuleName() string {
	return self.Request.GetRuleName()
}

// 获取请求中临时缓存数据
// defaultValue 不能为 interface{}(nil)
func (self *Context) GetTemp(key string, defaultValue interface{}) interface{} {
	return self.Request.GetTemp(key, defaultValue)
}

// 获取请求中全部缓存数据
func (self *Context) GetTemps() request.Temp {
	return self.Request.GetTemps()
}

// 获得一个请求的缓存数据副本。
func (self *Context) CopyTemps() request.Temp {
	temps := make(request.Temp)
	for k, v := range self.Request.GetTemps() {
		temps[k] = v
	}
	return temps
}

// 从原始请求获取Url，从而保证请求前后的Url完全相等，且中文未被编码。
func (self *Context) GetUrl() string {
	return self.Request.Url
}

func (self *Context) GetMethod() string {
	return self.Request.GetMethod()
}

func (self *Context) GetHost() string {
	return self.Response.Request.URL.Host
}

// 获取响应头信息。
func (self *Context) GetHeader() http.Header {
	return self.Response.Header
}

// 获取请求头信息。
func (self *Context) GetRequestHeader() http.Header {
	return self.Response.Request.Header
}

func (self *Context) GetReferer() string {
	return self.Response.Request.Header.Get("Referer")
}

// 获取响应的Cookie。
func (self *Context) GetCookie() string {
	return self.Response.Header.Get("Set-Cookie")
}

// GetHtmlParser returns goquery object binded to target crawl result.
func (self *Context) GetDom() *goquery.Document {
	if self.dom == nil {
		self.initDom()
	}
	return self.dom
}

// GetBodyStr returns plain string crawled.
func (self *Context) GetText() string {
	if self.text == nil {
		self.initText()
	}
	return util.Bytes2String(self.text)
}

//**************************************** 私有方法 *******************************************\\

// 获取规则。
func (self *Context) getRule(ruleName ...string) (name string, rule *Rule, found bool) {
	if len(ruleName) == 0 {
		if self.Response == nil {
			return
		}
		name = self.GetRuleName()
	} else {
		name = ruleName[0]
	}
	rule, found = self.spider.GetRule(name)
	return
}

// GetHtmlParser returns goquery object binded to target crawl result.
func (self *Context) initDom() *goquery.Document {
	if self.text == nil {
		self.initText()
	}
	var err error
	self.dom, err = goquery.NewDocumentFromReader(bytes.NewReader(self.text))
	if err != nil {
		return nil
	}
	return self.dom
}

// GetBodyStr returns plain string crawled.
func (self *Context) initText() {
	var err error

	// 采用surf内核下载时，尝试自动转码
	if self.Request.DownloaderID == request.SURF_ID {
		var contentType, pageEncode string
		// 优先从响应头读取编码类型
		contentType = self.Response.Header.Get("Content-Type")
		if _, params, err := mime.ParseMediaType(contentType); err == nil {
			if cs, ok := params["charset"]; ok {
				pageEncode = strings.ToLower(strings.TrimSpace(cs))
			}
		}
		// 响应头未指定编码类型时，从请求头读取
		if len(pageEncode) == 0 {
			contentType = self.Request.Header.Get("Content-Type")
			if _, params, err := mime.ParseMediaType(contentType); err == nil {
				if cs, ok := params["charset"]; ok {
					pageEncode = strings.ToLower(strings.TrimSpace(cs))
				}
			}
		}

		switch pageEncode {
		// 不做转码处理
		case "utf8", "utf-8", "unicode-1-1-utf-8":
		default:
			// 指定了编码类型，但不是utf8时，自动转码为utf8
			// get converter to utf-8
			// Charset auto determine. Use golang.org/x/net/html/charset. Get response body and change it to utf-8
			var destReader io.Reader

			if len(pageEncode) == 0 {
				destReader, err = charset.NewReader(self.Response.Body, "")
			} else {
				destReader, err = charset.NewReaderLabel(pageEncode, self.Response.Body)
			}

			if err == nil {
				self.text, err = ioutil.ReadAll(destReader)
				if err == nil {
					self.Response.Body.Close()
					return
				} else {
					logs.Log.Warning(" *     [convert][%v]: %v (ignore transcoding)", self.GetUrl(), err)
				}
			} else {
				logs.Log.Warning(" *     [convert][%v]: %v (ignore transcoding)", self.GetUrl(), err)
			}
		}
	}

	// 不做转码处理
	self.text, err = ioutil.ReadAll(self.Response.Body)
	self.Response.Body.Close()
	if err != nil {
		return
	}
}

// 生成文本结果。
// 用ruleName指定匹配的ItemFields字段，为空时默认当前规则。
func (self *Context) CreatItem(item map[int]interface{}, ruleName ...string) map[string]interface{} {
	_, rule, found := self.getRule(ruleName...)
	if !found {
		logs.Log.Error("蜘蛛 %s 调用CreatItem()时，指定的规则名不存在！", self.spider.GetName())
		return nil
	}

	var item2 = make(map[string]interface{}, len(item))
	for k, v := range item {
		field := self.spider.GetItemField(rule, k)
		item2[field] = v
	}
	return item2
}

// 输出文本结果。
// item类型为map[int]interface{}时，根据ruleName现有的ItemFields字段进行输出，
// item类型为map[string]interface{}时，ruleName不存在的ItemFields字段将被自动添加，
// ruleName为空时默认当前规则。
func (self *Context) Output(item interface{}, ruleName ...string) {
	_ruleName, rule, found := self.getRule(ruleName...)
	if !found {
		logs.Log.Error("spider: %s Output() error! rule name not exists！", self.spider.GetName())
		return
	}
	var _item map[string]interface{}
	switch item2 := item.(type) {
	case map[int]interface{}:
		_item = self.CreatItem(item2, _ruleName)
	case request.Temp:
		for k := range item2 {
			self.spider.UpsertItemField(rule, k)
		}
		_item = item2
	case map[string]interface{}:
		for k := range item2 {
			self.spider.UpsertItemField(rule, k)
		}
		_item = item2
	}
	self.Lock()
	self.items = append(self.items, data.GetDataCell(_ruleName, _item, self.GetUrl(), self.GetReferer(), time.Now().Format("2006-01-02 15:04:05")))
	self.Unlock()
}
