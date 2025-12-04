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

func firstLine(s string) string {
    if i := strings.IndexByte(s, '\n'); i >= 0 {
        return s[:i]
    }
    return s
}

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

