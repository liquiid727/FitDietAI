// 文件功能：Markdown 文本清洗与路径处理的辅助方法集合。
// 包功能：为具体解析器实现提供通用工具函数。
package impl

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	mdExtRegex      = regexp.MustCompile(`(?i)\.md$`)
	headerRegex     = regexp.MustCompile(`(?m)^#{1,6}\s+.*$`)
	imageMdRegex    = regexp.MustCompile(`!\[[^\]]*\]\([^\)]*\)`)
	htmlTagRegex    = regexp.MustCompile(`<[^>]+>`)
	multiBlankRegex = regexp.MustCompile(`\n{3,}`)
)

// cleanMarkdown：清洗 Markdown 文本，去除图片标记与 HTML 标签，归一化换行并压缩空行。
// 参数：
//   - s：原始 Markdown 文本。
//
// 返回：
//   - string：清洗后的文本。
func cleanMarkdown(s string) string {
	s = imageMdRegex.ReplaceAllString(s, "")
	s = htmlTagRegex.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.TrimSpace(s)
	s = multiBlankRegex.ReplaceAllString(s, "\n\n")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

// firstLine：返回文本的首行内容；若无换行符则返回全文。
// 参数：
//   - s：输入文本。
//
// 返回：
//   - string：首行文本。
func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// extractCategoryName：从相对路径抽取一级目录分类与文件基名。
// 参数：
//   - rel：相对路径。
//
// 返回：
//   - string：分类；
//   - string：文件名（不含扩展名）。
func extractCategoryName(rel string) (string, string) {
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) >= 2 {
		return parts[0], strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(parts[len(parts)-1]))
	}
	if len(parts) == 1 {
		return "", strings.TrimSuffix(parts[0], filepath.Ext(parts[0]))
	}
	return "", ""
}
