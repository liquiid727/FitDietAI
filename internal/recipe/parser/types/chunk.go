package types

// Chunk represents a parsed document segment with metadata.
type Chunk struct {
    ID       string `json:"id"`
    DocID    string `json:"doc_id"`
    Index    int    `json:"index"`
    Header   string `json:"header"`
    Text     string `json:"text"`
    Source   string `json:"source"`
    Category string `json:"category"`
    Name     string `json:"name"`
    Path     string `json:"path"`
}

// Options controls parsing behaviors.
type Options struct {
    ByHeader  bool
    ChunkSize int
    Overlap   int
    Timestamp bool
}

