// 文件功能：命令行入口，读取菜谱 Markdown 文档并委派内部解析器生成 JSONL 输出；提供基础 CLI 参数配置与文件写出服务。
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	parser "cook/internal/recipe/parser"
	"cook/internal/recipe/parser/impl"
	"cook/internal/recipe/parser/types"
)

// keep CLI flags and behavior; parsing is delegated to internal parser implementation

// main：解析 CLI 参数并调用内部解析器完成 Markdown→Chunk→JSONL 的生成；对错误进行标准输出与退出码处理。
func main() {
	var (
		dir       string
		out       string
		chunkSize int
		overlap   int
		byHeader  bool
	)
	flag.StringVar(&dir, "dir", defaultRecipesDir(), "recipes root directory")
	flag.StringVar(&out, "out", filepath.Join("parse", "out", "chunks.jsonl"), "output JSONL file path")
	flag.IntVar(&chunkSize, "chunk", 1200, "max characters per chunk (when not splitting by header)")
	flag.IntVar(&overlap, "overlap", 100, "overlap characters between chunks")
	flag.BoolVar(&byHeader, "byHeader", true, "split chunks using Markdown headers when possible")
	flag.Parse()

	var p parser.Parser = impl.NewMarkdownParser()
	files, err := p.Collect(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect files error: %v\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "no markdown files under: %s\n", dir)
		os.Exit(2)
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create out dir error: %v\n", err)
		os.Exit(3)
	}
	f, err := os.Create(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create out file error: %v\n", err)
		os.Exit(4)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	chunks, err := p.ParseFiles(files, types.Options{ByHeader: byHeader, ChunkSize: chunkSize, Overlap: overlap, Timestamp: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse files error: %v\n", err)
		os.Exit(5)
	}
	for _, c := range chunks {
		if err := writeJSONL(w, c); err != nil {
			fmt.Fprintf(os.Stderr, "write jsonl error: %v\n", err)
		}
	}
}

// defaultRecipesDir：自动探测默认菜谱目录；优先使用 FitDietAI/recipes，其次兼容 recipes/recipies。
func defaultRecipesDir() string {
	candidates := []string{
		filepath.Join("FitDietAI", "recipes"),
		"recipes",
		filepath.Join("FitDietAI", "recipies"),
		"recipies",
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c
		}
	}
	return filepath.Join("FitDietAI", "recipes")
}

// writeJSONL：将单个分块以一行 JSON 形式写入到输出文件；调用方负责缓冲区刷新与文件关闭。
func writeJSONL(w *bufio.Writer, c types.Chunk) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	if _, err := w.WriteString("\n"); err != nil {
		return err
	}
	return nil
}
