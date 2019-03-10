package agents

import "github.com/TeaWeb/code/teaconfigs/forms"

// 数据源接口定义
type SourceInterface interface {
	// 名称
	Name() string

	// 代号
	Code() string

	// 描述
	Description() string

	// 校验
	Validate() error

	// 执行
	Execute(params map[string]string) (value interface{}, err error)

	// 获得数据格式
	DataFormatCode() SourceDataFormat

	// 表单信息
	Form() *forms.Form

	// 显示信息
	Presentation() *forms.Presentation
}
