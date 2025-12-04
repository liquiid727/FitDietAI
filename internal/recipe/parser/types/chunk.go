// 文件功能：Chunk 分块结构与 Options 解析选项定义。
// 包功能：为解析器实现与调用方提供统一的数据结构接口。
package types

// Chunk represents a parsed document segment with metadata.
// Chunk：解析得到的文档分块。
//   - ID：分块唯一标识；
//   - DocID：所属文档标识；
//   - Index：分块在文档内的顺序索引；
//   - Header：分块标题（若为按长度分块则为占位标题）；
//   - Text：分块正文内容；
//   - Source：来源路径与可选时间戳；
//   - Category：文档一级目录分类；
//   - Name：文档基名（不含后缀）；
//   - Path：文档相对路径。
type Chunk struct {
	ID       string `json:"id"`       // 分块唯一标识
	DocID    string `json:"doc_id"`   // 文档标识
	Index    int    `json:"index"`    // 分块序号
	Header   string `json:"header"`   // 分块标题
	Text     string `json:"text"`     // 分块正文
	Source   string `json:"source"`   // 来源信息
	Category string `json:"category"` // 分类
	Name     string `json:"name"`     // 文档名
	Path     string `json:"path"`     // 相对路径
}

// Options controls parsing behaviors.
// Options：解析行为配置。
//   - ByHeader：是否按标题分块；
//   - ChunkSize：按长度分块的最大字符数；
//   - Overlap：相邻分块之间的重叠字符数；
//   - Timestamp：是否在 Source 附加 UTC 时间戳。
type Options struct {
	ByHeader  bool // 标题分块开关
	ChunkSize int  // 分块最大长度
	Overlap   int  // 分块重叠长度
	Timestamp bool // Source 是否带时间戳
}
