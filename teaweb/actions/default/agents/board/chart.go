package apps

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teadb"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/board/scripts"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ChartAction actions.Action

func (this *ChartAction) RunGet(params struct {
	AgentId string
	AppId   string
	ItemId  string
	ChartId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	chart := item.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("找不到Board")
	}

	boardChart := board.FindChart(params.ChartId)
	if boardChart == nil {
		this.Fail("找不到BoardChart")
	}

	if len(boardChart.Name) == 0 {
		boardChart.Name = chart.Name
	}

	this.Data["timePasts"] = teaconfigs.AllTimePasts()

	if len(boardChart.TimeType) == 0 {
		boardChart.TimeType = "default"
	}
	if len(boardChart.TimePast) == 0 {
		boardChart.TimePast = teaconfigs.TimePast1h
	}

	if len(boardChart.DayTo) == 0 {
		boardChart.DayTo = timeutil.Format("Y-m-d")
	}

	this.Data["agentId"] = params.AgentId
	this.Data["appId"] = params.AppId
	this.Data["itemId"] = params.ItemId
	this.Data["chartId"] = params.ChartId

	this.Data["chart"] = maps.Map{
		"name":     boardChart.Name,
		"timeType": boardChart.TimeType,
		"timePast": boardChart.TimePast,
		"dayFrom":  boardChart.DayFrom,
		"dayTo":    boardChart.DayTo,
	}

	this.Show()
}

// 获取数据
func (this *ChartAction) RunPost(params struct {
	Name     string
	AgentId  string
	AppId    string
	ItemId   string
	ChartId  string
	TimeType string
	TimePast string
	DayFrom  string
	DayTo    string

	Must *actions.Must
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	chart := item.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("找不到Board")
	}

	boardChart := board.FindChart(params.ChartId)
	if boardChart == nil {
		this.Fail("找不到BoardChart")
	}

	o, err := chart.AsObject()
	if err != nil {
		this.Fail("数据错误：" + err.Error())
	}

	code, err := o.AsJavascript(maps.Map{
		"name":    params.Name,
		"columns": chart.Columns,
	})
	if err != nil {
		this.Fail("数据错误：" + err.Error())
	}

	mongoEnabled := teadb.SharedDB().Test() == nil
	engine := scripts.NewEngine()
	engine.SetMongo(mongoEnabled)
	engine.SetCache(false)

	engine.SetContext(&scripts.Context{
		Agent:    agent,
		App:      app,
		Item:     item,
		TimeType: params.TimeType,
		TimePast: params.TimePast,
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
	})

	widgetCode := `var widget = new widgets.Widget({
	"name": "看板",
	"requirements": ["mongo"]
});

widget.run = function () {
`
	widgetCode += "{\n" + code + "\n}\n"
	widgetCode += `
};
`

	err = engine.RunCode(widgetCode)
	if err != nil {
		this.Fail("发生错误：" + err.Error())
	}

	this.Data["charts" ] = engine.Charts()
	this.Data["output"] = engine.Output()
	this.Success()
}
