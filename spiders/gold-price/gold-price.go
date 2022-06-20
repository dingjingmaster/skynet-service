package gold_price

import (
	"skynet-service/app/base/spider"
	"skynet-service/app/common/request"
)

func init() {
	GoldPrice.Register()
}

type GoldPriceData struct {
	Ts int64 `json:"ts"`
	Tsj int64 `json:"tsj"`
	Date string `json:"date"`
	Items [] struct {
		Curr string `json:"curr"`
		XauPrice float64 `json:"xauPrice"`
		XagPrice float64 `json:"xagPrice"`
		ChgXau float64 `json:"chgXau"`
		ChgXag float64 `json:"chgXag"`
		PcXau float64 `json:"pcXau"`
		PcXag float64 `json:"pcXag"`
		XauClose float64 `json:"xauClose"`
		XagClose float64 `json:"xagClose"`
	} `json:"items"`
}

var GoldPrice = &spider.Spider{
	Name:        "金价抓取",
	Description: "获取伦敦金价、美国金价、中国金价",
	RequestType: request.REQUEST_TYPE_PUT,
	RuleTree:    nil,
}