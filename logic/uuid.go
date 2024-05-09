package logic

import (
	"VideoWeb/Utilities/logf"
	"VideoWeb/define"
	"github.com/yitter/idgenerator-go/idgen"
)

func init() {
	var options = idgen.NewIdGeneratorOptions(1)
	options.BaseTime = define.BaseTime
	idgen.SetIdGenerator(options)
	logf.WriteInfoLog("logic::uuid::init", "初始化雪花算法完成!")
}

// GetUUID 生成UUID
func GetUUID() int64 {
	return idgen.NextId()
}
