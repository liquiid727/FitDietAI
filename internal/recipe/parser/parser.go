package parser

import (
    "bytes"
    "io/fs"
    "os"
    "path/filepath"
    "github.com/yuin/goldmark"
    meta "github.com/yuin/goldmark-meta"
)

type Recipe struct {
    Title        string
    Servings     int
    CookingTime  string
    Difficulty   string
    Tags         []string
    Ingredients  []string
    Steps        []string
    Notes        string
    RawMarkdown  string
}

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