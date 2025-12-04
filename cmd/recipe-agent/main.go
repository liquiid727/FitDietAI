// 文件功能：命令行入口；启动 Recipe Agent 的 CLI。
package main

import (
	"cook/internal/recipe/server"
	"log"
)

// main：执行 CLI 根命令；失败时以致命日志退出。
func main() {
	if err := server.Execute(); err != nil {
		log.Fatal(err)
	}
}
