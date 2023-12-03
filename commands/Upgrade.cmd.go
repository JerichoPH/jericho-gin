package commands

import (
	"jericho-gin/database"
)

// UpgradeCommand 程序升级
type UpgradeCommand struct{}

// NewUpgradeCommand 构造函数
func NewUpgradeCommand() UpgradeCommand {
	return UpgradeCommand{}
}

// 执行数据库管理语句
func (UpgradeCommand) execStatements(sql []string) {
	if len(sql) > 0 {
		for _, s := range sql {
			database.NewGormLauncher().GetConn("").Exec(s)
		}
	}
}

// Do 执行命令
func (receiver UpgradeCommand) Do(params []string) []string {
	switch params[0] {
	}

	return []string{"执行完成"}
}
