package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teautils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 路径规则详情
func (this *DetailAction) Run(params struct {
	ServerId   string
	LocationId string
}) {
	server, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "detail")

	if len(location.Pages) == 0 {
		location.Pages = []*teaconfigs.PageConfig{}
	}

	this.Data["location"] = maps.Map{
		"on":              location.On,
		"id":              location.Id,
		"type":            location.PatternType(),
		"pattern":         location.PatternString(),
		"name":            location.Name,
		"isBreak":         location.IsBreak,
		"caseInsensitive": location.IsCaseInsensitive(),
		"reverse":         location.IsReverse(),
		"root":            location.Root,
		"charset":         location.Charset,
		"index":           location.Index,
		"urlPrefix":       location.URLPrefix,
		"maxBodySize":     location.MaxBodySize,
		"enableStat":      !location.DisableStat,
		"gzip":            location.Gzip,
		"redirectToHttps": location.RedirectToHttps,
		"conds":           location.Cond,

		"fastcgi":     location.Fastcgi,
		"headers":     location.Headers,
		"cachePolicy": location.CachePolicy,
		"rewrite":     location.Rewrite,
		"websocket":   location.Websocket,
		"backends":    location.Backends,
		"wafId":       location.WafId,

		"shutdown": location.Shutdown,
		"pages":    location.Pages,
	}
	this.Data["server"] = server

	// 字符集
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets

	this.Data["accessLogs"] = proxyutils.FormatAccessLog(location.AccessLog)

	this.Show()
}
