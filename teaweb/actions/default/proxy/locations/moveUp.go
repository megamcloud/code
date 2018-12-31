package locations

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type MoveUpAction actions.Action

func (this *MoveUpAction) Run(params struct {
	Filename string
	Index    int
}) {
	proxy, err := teaconfigs.NewServerConfigFromFile(params.Filename)
	if err != nil {
		this.Fail(err.Error())
	}

	if params.Index >= 1 && params.Index < len(proxy.Locations) {
		prev := proxy.Locations[params.Index-1]
		current := proxy.Locations[params.Index]
		proxy.Locations[params.Index-1] = current
		proxy.Locations[params.Index] = prev
	}

	proxy.Save()

	proxyutils.NotifyChange()

	this.Refresh().Success()
}
