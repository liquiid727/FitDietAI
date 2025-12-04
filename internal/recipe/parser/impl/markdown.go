// 文件功能：Markdown 文档解析为结构化 Chunk（分块）并输出元数据；实现 Parser 接口的具体 Markdown 解析器。
// 创建日期：2025-12-04；最后修改日期：2025-12-04。
package impl

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cook/internal/recipe/parser/types"
)

// MarkdownParser：将 Markdown 文件解析为若干 Chunk；支持按标题（Header）与按长度（Size）两种分块策略。
// 业务背景：用于 RAG 数据准备阶段，将菜谱 Markdown 规范化、清洗并切分，以便后续 Embedding 与向量索引构建。
type MarkdownParser struct{}

// NewMarkdownParser：构造 MarkdownParser 实例。
// 功能说明：返回一个无状态的解析器；方法可并发使用，但调用方需自行控制并发与 I/O。
// 参数说明：无。
// 返回值说明：返回实现 Parser 的具体 Markdown 解析器。
// 示例代码：
//

func NewMarkdownParser() *MarkdownParser { return &MarkdownParser{} }

// Collect：在指定根目录下收集全部 Markdown 文件路径。
// 功能说明：递归遍历 root 目录；仅收集扩展名为 .md 的文件；对不可达或非法目录进行错误上抛（error wrapping）。
// 参数说明：
//   - root：扫描的根目录路径；必须为存在且可读的目录；空字符串将报错。
//
// 返回值说明：
//   - []string：匹配到的 Markdown 文件的绝对或相对路径集合；不保证排序稳定。
//   - error：当 root 为空、非目录或遍历期间发生 I/O 错误时返回具体错误。
func (p *MarkdownParser) Collect(root string) ([]string, error) {
	if root == "" {
		return nil, fmt.Errorf("collect: root is empty")
	}
	st, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("collect: stat root: %w", err)
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("collect: root not directory: %s", root)
	}
	out := make([]string, 0, 256)
	// 关键逻辑：使用 WalkDir 高效遍历目录树；过滤目录项；仅匹配 .md 扩展名（大小写不敏感）。
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		if mdExtRegex.MatchString(strings.ToLower(filepath.Ext(p))) {
			out = append(out, p)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect: walk: %w", err)
	}
	return out, nil
}

// ParseFiles：解析文件集合并输出 Chunk 列表。
// 功能说明：对传入的 Markdown 文件进行统一清洗（去图片标记、去 HTML 标签、归一化换行、压缩空行），并根据 Options 选择分块策略；
//
//	同时填充来源元数据（source）、分类（category）、文件名（name）等，保证后续检索与回溯能力。
//
// 参数说明：
//   - paths：待解析的 Markdown 文件列表；为空时返回错误。
//   - opts：解析选项，包含 ByHeader／ChunkSize／Overlap／Timestamp；其中：
//   - ByHeader：是否按标题分块；为 false 时使用按长度分块；
//   - ChunkSize：长度分块的最大字符数；≤0 时使用默认 1200；
//   - Overlap：长度分块的重叠字符数；<0 视为 0；
//   - Timestamp：source 字段是否附加 UTC 时间戳。
//
// 返回值说明：
//   - []types.Chunk：解析得到的分块列表；每个分块具备唯一 ID、索引、正文与元数据。
//   - error：任何 I/O 读取错误与参数非法将返回错误；默认使用 %w 包裹以便调用方诊断。
//
// 示例代码：
//
//	```go
//	p := impl.NewMarkdownParser()
//	files, _ := p.Collect("recipes")
//	opts := types.Options{ByHeader:false, ChunkSize:800, Overlap:100, Timestamp:true}
//	chunks, err := p.ParseFiles(files, opts)
//	if err != nil { /* 处理错误 */ }
//	```
func (p *MarkdownParser) ParseFiles(paths []string, opts types.Options) ([]types.Chunk, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("parse: no input files")
	}
	// defaults
	if opts.ChunkSize <= 0 {
		opts.ChunkSize = 1200
	}
	if opts.Overlap < 0 {
		opts.Overlap = 0
	}
	if !opts.ByHeader && opts.ChunkSize < 200 {
		opts.ChunkSize = 200
	}
	chunks := make([]types.Chunk, 0, len(paths)*4)
	docCounter := 0
	chunkCounter := 0
	for _, pth := range paths {
		// 关键逻辑：逐文件读取内容；对读取失败直接中止并返回错误，保证数据一致性；
		b, err := os.ReadFile(pth)
		if err != nil {
			return nil, fmt.Errorf("parse: read %s: %w", pth, err)
		}
		rel := pth
		// attempt to make path relative to root of first input
		// 关键逻辑：相对路径标准化；以首个输入的父目录作为“公共根”，便于生成稳定的 source/path 元数据；
		if root := commonRoot(paths); root != "" {
			if absRoot, e := filepath.Abs(root); e == nil {
				if absFile, e2 := filepath.Abs(pth); e2 == nil {
					if strings.HasPrefix(absFile, absRoot) {
						rel = strings.TrimPrefix(absFile, absRoot)
						rel = strings.TrimPrefix(rel, string(filepath.Separator))
					}
				}
			}
		}
		cat, name := extractCategoryName(rel)
		docID := fmt.Sprintf("doc-%d", docCounter)
		docCounter++
		// 关键逻辑：统一清洗 Markdown 文本；剔除图片标记与 HTML 标签，统一换行并压缩空行；
		text := cleanMarkdown(string(b))
		var docChunks []types.Chunk
		if opts.ByHeader {
			// 算法说明（标题分块）：
			//  1）按正则 ^#{1,6}\s+ 捕获所有标题位置；
			//  2）每个标题到下一个标题的区间作为一个分块；
			//  3）分块的 Header 取首行标题文本；Text 取标题后的正文；
			docChunks = splitByHeaders(text, docID, rel, cat, name, opts)
		} else {
			// 算法说明（长度分块）：
			//  1）对清洗后的文本按字符长度进行滑窗切分；
			//  2）相邻分块之间可配置 Overlap；
			//  3）Header 使用占位标记以体现块序号；
			docChunks = splitBySize(text, docID, rel, cat, name, opts)
		}
		for i := range docChunks {
			chunkCounter++
			docChunks[i].ID = fmt.Sprintf("chunk-%d", chunkCounter)
			docChunks[i].Index = i
			chunks = append(chunks, docChunks[i])
		}
	}
	return chunks, nil
}

func splitByHeaders(text, docID, rel, cat, name string, opts types.Options) []types.Chunk {
	idxs := headerRegex.FindAllStringIndex(text, -1)
	if len(idxs) == 0 {
		return []types.Chunk{{
			DocID:    docID,
			Header:   "",
			Text:     text,
			Source:   source(rel, opts.Timestamp),
			Category: cat,
			Name:     name,
			Path:     rel,
		}}
	}
	out := make([]types.Chunk, 0, len(idxs)+1)
	for i := 0; i < len(idxs); i++ {
		start := idxs[i][0]
		end := len(text)
		if i+1 < len(idxs) {
			end = idxs[i+1][0]
		}
		seg := strings.TrimSpace(text[start:end])
		headerLine := firstLine(seg)
		body := strings.TrimSpace(strings.TrimPrefix(seg, headerLine))
		out = append(out, types.Chunk{
			DocID:    docID,
			Header:   strings.TrimSpace(headerLine),
			Text:     body,
			Source:   source(rel, opts.Timestamp),
			Category: cat,
			Name:     name,
			Path:     rel,
		})
	}
	return out
}

func splitBySize(text, docID, rel, cat, name string, opts types.Options) []types.Chunk {
	size := opts.ChunkSize
	overlap := opts.Overlap
	out := make([]types.Chunk, 0, strings.Count(text, "\n")/20+1)
	runes := []rune(text)
	start := 0
	idx := 0
	for start < len(runes) {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}
		body := strings.TrimSpace(string(runes[start:end]))
		out = append(out, types.Chunk{
			DocID:    docID,
			Header:   fmt.Sprintf("# chunk %d", idx),
			Text:     body,
			Source:   source(rel, opts.Timestamp),
			Category: cat,
			Name:     name,
			Path:     rel,
		})
		if end == len(runes) {
			break
		}
		start = end - overlap
		if start < 0 {
			start = 0
		}
		idx++
	}
	return out
}

func source(rel string, ts bool) string {
	if ts {
		return fmt.Sprintf("%s|%s", rel, time.Now().UTC().Format(time.RFC3339))
	}
	return rel
}

func commonRoot(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	dir := filepath.Dir(paths[0])
	return dir
}
