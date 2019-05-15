package teaconfigs

import (
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/logs"
)

// 节点配置文件名
var nodeConfigFile = "node.conf"

// 节点配置
type NodeConfig struct {
	Id            string   `yaml:"id" json:"id"`                       // ID
	On            bool     `yaml:"on" json:"on"`                       // 是否启用
	Name          string   `yaml:"name" json:"name"`                   // 名称
	ClusterId     string   `yaml:"clusterId" json:"clusterId"`         // 集群ID
	ClusterSecret string   `yaml:"clusterSecret" json:"clusterSecret"` // 集群秘钥
	ClusterAddr   string   `yaml:"clusterAddr" json:"clusterAddr"`     // 集群通讯地址
	Role          NodeRole `yaml:"role" json:"role"`                   // 角色
}

// 取得当前节点配置
func SharedNodeConfig() *NodeConfig {
	configFile := files.NewFile(Tea.ConfigFile(nodeConfigFile))
	if !configFile.Exists() {
		return nil
	}

	data, err := configFile.ReadAll()
	if err != nil {
		logs.Error(err)
		return nil
	}

	config := &NodeConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil
	}

	return config
}

// 保存到文件
func (this *NodeConfig) Save() error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	configFile := files.NewFile(Tea.ConfigFile(nodeConfigFile))
	return configFile.Write(data)
}

// 是否为Master
func (this *NodeConfig) IsMaster() bool {
	return this.Role == NodeRoleMaster
}