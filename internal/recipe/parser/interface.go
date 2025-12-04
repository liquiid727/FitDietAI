// Package parser defines core interfaces for recipe markdown parsing.
package parser

import "cook/internal/recipe/parser/types"

// Parser is the core interface to parse markdown files into chunks.
type Parser interface {
    Collect(root string) ([]string, error)
    ParseFiles(paths []string, opts types.Options) ([]types.Chunk, error)
}

