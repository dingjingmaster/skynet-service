package gold_price

import (
	"encoding/json"
	"fmt"
	"skynet-service/app/common/timeUtils"
	"skynet-service/app/logs"
	"skynet-service/app/spider"

	//"skynet-service/app/common/request"
	"skynet-service/app/downloader/request"
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
	RuleTree:    &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.AddQueue(&request.Request{
				Url:  "https://data-asg.goldprice.org/dbXRates/CNY",
				Rule: "获取标准黄金价格",
				Temp: map[string]interface{}{
					"target": "first",
				},
			})
		},

		Trunk: map[string]*spider.Rule{
			"获取标准黄金价格": {
				ItemFields: []string {
					"gid",
					"gold-1",
					"gold-2",
					"gold-h",
					"gold-l",
					"gold-v",
					"gold-a",
					"gold-t",
				},
				ParseFunc: func(ctx *spider.Context) {
					jsonStr := ctx.GetText()
					gs := &GoldPriceData{}
					json.Unmarshal([]byte(jsonStr), &gs)
					gold1 := gs.Items[0].XauPrice
					gold2 := gs.Items[0].XauPrice
					goldh := gs.Items[0].XauPrice
					goldl := gs.Items[0].XauPrice
					goldv := gs.Items[0].ChgXau
					golda := 0
					goldt := timeUtils.GetTimeStamp()
					id := fmt.Sprintf("%d-%.2f-%.2f-%.2f-%.2f-%.2f-%d", goldt, gold1, gold2, goldh, goldl, goldv, golda)

					ctx.Output (map[int]interface{} {
						0: id,
						1: gold1,
						2: gold2,
						3: goldh,
						4: goldl,
						5: goldv,
						6: golda,
						7: goldt,
					})
					logs.Log.Informational("%v", id)
				},
			},
		},
	},
}