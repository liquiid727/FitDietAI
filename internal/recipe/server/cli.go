// 文件功能：命令行入口与子命令定义；提供服务启动与索引构建指令。
// 包功能：server 包，封装 CLI 与 HTTP 服务相关逻辑。
package server

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "recipe-agent",
	Short: "Recipe document assistant Agent",
}

// Execute：执行根命令。
// 功能说明：挂载子命令并处理命令行输入；返回可能的执行错误。
// 参数说明：无。
// 返回值说明：
//   - error：命令执行失败的错误。
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(indexCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start REST server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startHTTP()
	},
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index recipes from recipes/",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(os.Stdout, "indexing recipes: TODO")
		return nil
	},
}
