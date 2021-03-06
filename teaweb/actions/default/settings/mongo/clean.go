package mongo

import (
	"github.com/TeaWeb/code/teaconfigs/db"
	"github.com/iwind/TeaGo/actions"
)

type CleanAction actions.Action

// 设置自动清理
func (this *CleanAction) Run(params struct{}) {
	config, _ := db.LoadMongoConfig()
	if config != nil {
		this.Data["accessLog"] = config.AccessLog
	} else {
		this.Data["accessLog"] = nil
	}

	this.Show()
}
