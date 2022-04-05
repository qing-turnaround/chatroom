package global

import (
	"os"
	"path/filepath"
)

// 项目主目录
var RootDir string

func init() {
	inferRootDir()
	// 初始化配置
	configInit()
	// 初始化zap
	zapInit()
}

// 推断 Root目录
func inferRootDir() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(string) string
	infer = func(dir string) string {
		if exists(dir + "/template") {
			return dir
		}

		// 查看dir的父目录
		parent := filepath.Dir(dir)
		return infer(parent)
	}

	RootDir = infer(pwd)
}

func exists(dir string) bool {
	// 查找主机是不是存在 dir
	_, err := os.Stat(dir)
	return err == nil || os.IsExist(err)
}
