// 文件功能：解析器核心接口定义；约定 Markdown 文档收集与解析为 Chunk 的统一能力。
// 包功能：parser 包，提供 Parser 接口，便于不同实现（如 Markdown、HTML 等）进行扩展与替换。
package parser

import "cook/internal/recipe/parser/types"

// Parser：解析器接口；用于将原始文档解析为结构化 Chunk。
//   - Collect：收集指定根目录下的文档路径集合；
//   - ParseFiles：按 Options 对文件进行解析与分块输出。
type Parser interface {
	Collect(root string) ([]string, error)
	ParseFiles(paths []string, opts types.Options) ([]types.Chunk, error)
}
