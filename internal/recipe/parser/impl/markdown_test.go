// 文件功能：Markdown 解析器的单元测试与基准测试；验证收集与分块功能的正确性与性能特征。
package impl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cook/internal/recipe/parser/types"
)

func TestCollectAndParseByHeader(t *testing.T) {
	dir := t.TempDir()
	md := "# 标题一\n内容A\n\n## 子标题\n内容B\n\n# 标题二\n内容C"
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte(md), 0o644); err != nil {
		t.Fatalf("write md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte(md), 0o644); err != nil { // should be ignored
		t.Fatalf("write txt: %v", err)
	}
	p := NewMarkdownParser()
	files, err := p.Collect(dir)
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expect 1 md file, got %d", len(files))
	}
	chunks, err := p.ParseFiles(files, types.Options{ByHeader: true, Timestamp: false})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(chunks) != 3 {
		t.Fatalf("expect 3 chunks by headers, got %d", len(chunks))
	}
	if !strings.HasPrefix(chunks[0].Header, "# ") {
		t.Fatalf("first chunk header not found: %q", chunks[0].Header)
	}
}

func TestParseBySize(t *testing.T) {
	dir := t.TempDir()
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString("行内容\n")
	}
	if err := os.WriteFile(filepath.Join(dir, "c.md"), []byte(sb.String()), 0o644); err != nil {
		t.Fatalf("write md: %v", err)
	}
	p := NewMarkdownParser()
	files, err := p.Collect(dir)
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	chunks, err := p.ParseFiles(files, types.Options{ByHeader: false, ChunkSize: 50, Overlap: 10, Timestamp: true})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(chunks) == 0 {
		t.Fatalf("expected non-empty chunks")
	}
	if !strings.Contains(chunks[0].Source, ".md") {
		t.Fatalf("source missing path")
	}
}

func BenchmarkCleanMarkdown(b *testing.B) {
	text := strings.Repeat("![](img) <b>x</b>  内容\n\n\n", 500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cleanMarkdown(text)
	}
}

func BenchmarkParseByHeader(b *testing.B) {
	dir := b.TempDir()
	var sb strings.Builder
	for i := 0; i < 500; i++ {
		sb.WriteString("# H\n内容\n")
	}
	_ = os.WriteFile(filepath.Join(dir, "d.md"), []byte(sb.String()), 0o644)
	p := NewMarkdownParser()
	files, _ := p.Collect(dir)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.ParseFiles(files, types.Options{ByHeader: true})
	}
}
