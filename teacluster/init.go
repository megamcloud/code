package teacluster

import (
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teahooks"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"time"
)

func init() {
	if !ClusterEnabled {
		return
	}

	// register actions
	RegisterActionType(
		new(SuccessAction),
		new(FailAction),
		new(RegisterAction),
		new(PushAction),
		new(PullAction),
		new(NotifyAction),
		new(SumAction),
		new(SyncAction),
		new(PingAction),
	)

	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// build
		SharedManager.BuildSum()

		// start manager
		go func() {
			ticker := time.NewTicker(60 * time.Second)
			for {
				err := SharedManager.Start()
				if err != nil {
					logs.Println("[cluster]" + err.Error())
				}

				// retry N seconds later
				select {
				case <-ticker.C:
					// every N seconds
				case <-SharedManager.RestartChan:
					// retry immediately
				}
			}
		}()

		// start ping
		timers.Loop(60*time.Second, func(looper *timers.Looper) {
			node := teaconfigs.SharedNodeConfig()
			if node != nil && node.On && SharedManager.IsActive() {
				err := SharedManager.Write(&PingAction{})
				if err != nil {
					logs.Println("[cluster]" + err.Error())
				}
			}
		})
	})

	TeaGo.BeforeStop(func(server *TeaGo.Server) {
		if SharedManager != nil {
			SharedManager.Stop()
		}
	})

	teahooks.On(teahooks.EventConfigChanged, func() {
		node := teaconfigs.SharedNodeConfig()
		if node != nil && node.On {
			SharedManager.SetIsChanged(true)
		}
	})
}
