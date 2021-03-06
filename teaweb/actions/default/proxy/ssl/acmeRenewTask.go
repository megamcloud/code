package ssl

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"time"
)

type AcmeRenewTaskAction actions.Action

// 删除ACME任务
func (this *AcmeRenewTaskAction) RunPost(params struct {
	ServerId string
	TaskId   string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.SSL == nil {
		this.Success()
	}

	task := server.SSL.FindCertTask(params.TaskId)
	if task == nil {
		this.Fail("找不到任务")
	}
	task.RunAt = time.Now().Unix()

	client, err := task.Request.Client()
	if err != nil {
		task.RunError = "认证失败：" + err.Error()
		err = server.Save()
		if err != nil {
			logs.Error(err)
		}
		this.Fail(task.RunError)
	}

	err = task.Request.Renew(client)
	if err != nil {
		task.RunError = "更新失败：" + err.Error()
		err = server.Save()
		if err != nil {
			logs.Error(err)
		}
		this.Fail(task.RunError)
	}

	task.RunError = ""
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 更新证书
	found := false
	for _, cert := range server.SSL.Certs {
		if cert.TaskId == task.Id {
			err = task.Request.WriteCertFile(Tea.ConfigFile(cert.CertFile))
			if err != nil {
				logs.Error(err)
			}

			err = task.Request.WriteKeyFile(Tea.ConfigFile(cert.KeyFile))
			if err != nil {
				logs.Error(err)
			}

			found = true
		}
	}

	// 通知更新
	if found {
		proxyutils.NotifyChange()
	}

	this.Success()
}
