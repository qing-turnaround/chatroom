package logic

import (
	"chatroom/global"

	"github.com/importcjj/sensitive"
)

// 过滤 敏感词
func FilterSensitive(content string) string {
	sensitiveString := global.Conf.Sensitive
	filter := sensitive.New()
	filter.AddWord(sensitiveString...)
	return filter.Replace(content, '*')
}