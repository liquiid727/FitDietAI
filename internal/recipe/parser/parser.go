// 文件功能：基础 Markdown 解析示例；提供原始 Recipe 结构与目录/文件解析方法。
// 包功能：parser 包，除接口外还包含早期版本的简易解析器。
package parser

import (
	"bytes"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"io/fs"
	"os"
	"path/filepath"
)

// Recipe：菜谱文档的基础结构体（示例）。
//   - Title：菜谱标题；
//   - Servings：份量；
//   - CookingTime：烹饪时长；
//   - Difficulty：难度评估；
//   - Tags：标签集合；
//   - Ingredients：原料清单；
//   - Steps：操作步骤；
//   - Notes：补充说明；
//   - RawMarkdown：原始 Markdown 文本。
type Recipe struct {
	Title       string   // 菜谱标题
	Servings    int      // 份量
	CookingTime string   // 烹饪时长
	Difficulty  string   // 难度
	Tags        []string // 标签
	Ingredients []string // 原料
	Steps       []string // 步骤
	Notes       string   // 备注
	RawMarkdown string   // 原始 Markdown
}

// ParseDir：解析目录下的所有 Markdown 文件为 Recipe 结构。
// 参数：
//   - dir：目录路径。
//
// 返回：
//   - []*Recipe：解析结果集合；
//   - error：遍历或解析过程中产生的错误。
func ParseDir(dir string) ([]*Recipe, error) {
	var out []*Recipe
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		r, e := ParseFile(path)
		if e != nil {
			return e
		}
		out = append(out, r)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ParseFile：解析单个 Markdown 文件为 Recipe；当前实现保留原始 Markdown 文本。
// 参数：
//   - path：文件路径。
//
// 返回：
//   - *Recipe：解析后的菜谱结构；
//   - error：读取或转换失败返回错误。
func ParseFile(path string) (*Recipe, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	md := goldmark.New(goldmark.WithExtensions(meta.New()))
	var buf bytes.Buffer
	if err := md.Convert(b, &buf); err != nil {
		return nil, err
	}
	return &Recipe{RawMarkdown: string(b)}, nil
}
