## 目标
- 初始化标准 Go 模块与目录结构，搭建可启动的服务端与命令行框架
- 选定并拉取基础依赖：Eino（模型/向量/检索）、HTTP/Swagger、配置、Markdown 解析
- 预置 DeepSeek 模型接入（OpenAI 兼容 BaseURL），为后续 RAG 问答打好地基

## 目录结构
- `cmd/recipe-agent/main.go` 主入口（启动 REST/CLI）
- `internal/recipe/`
  - `parser/` Markdown 解析与食谱元数据抽取
  - `vector/` 索引与检索封装（Eino Indexer/Retriever/Embedding）
  - `qa/` 问答流水线与对话历史管理（Eino ChatModel + RAG）
  - `server/` REST 路由与处理器（Chi + Swagger）
  - `config/` 配置加载（Viper），敏感信息走环境变量
  - `store/` 运行期存储（如内存会话、缓存等）
- `pkg/` 可复用通用工具
- `configs/` 配置文件（如 `config.yaml`），不含密钥
- `recipes/` 示例食谱 Markdown（后续用于索引）

## 基础依赖与版本
- Eino 核心与扩展
  - `github.com/cloudwego/eino@latest`
  - `github.com/cloudwego/eino-ext/components/model/openai@latest`（将 BaseURL 指向 DeepSeek）
  - `github.com/cloudwego/eino-ext/components/embedding/openai@latest`（用于向量化；如需本地/其他嵌入可替换）
  - `github.com/cloudwego/eino-ext/components/retriever/es8@latest`（Elasticsearch 8 向量检索；也可后续替换 VikingDB 等）
- Web/Swagger
  - `github.com/go-chi/chi/v5@latest`
  - `github.com/swaggo/swag@latest`
  - `github.com/swaggo/http-swagger@latest`
  - `github.com/swaggo/files@latest`
- 配置与 CLI
  - `github.com/spf13/viper@latest`
  - `github.com/spf13/cobra@latest`
- Markdown
  - `github.com/yuin/goldmark@latest`
  - `github.com/yuin/goldmark-meta@latest`

说明与依据：
- Eino 官方：用户手册与组件结构 [5]、组件索引器 [2]、检索器 [4]、ES8 检索实现 [3]
- DeepSeek API：OpenAI 兼容 BaseURL `https://api.deepseek.com` 或 `https://api.deepseek.com/v1` [1]

## 初始化步骤（获批后执行的具体命令）
- 创建模块
  - `go mod init github.com/your-org/recipe-agent`
  - `go get` 拉取上述依赖
- 目录与骨架
  - 建立上述目录；创建 `main.go`、`internal/...` 基础文件
  - `configs/config.yaml`（示例：`server.port`、`retriever.es8.*`、`model.deepseek.*` 等）
- 配置加载
  - Viper 读取 `configs/config.yaml`，并允许环境变量覆盖（如 `DEEPSEEK_API_KEY`、`OPENAI_API_KEY`、`ES_ADDRESS`）
- HTTP 服务
  - Chi 路由：`GET /health`、`GET /api/v1/recipes`（占位返回）、`POST /api/v1/query`（占位——后续接 RAG）
  - 集成 Swagger：生成 OpenAPI 文档并通过 `GET /swagger/index.html` 提供
- CLI
  - `recipe-agent index`：扫描 `recipes/`，解析并索引到向量库
  - `recipe-agent serve`：启动 REST 服务

## 核心模块设计（第一阶段落地范围）
- Markdown 解析器
  - 采用 Front Matter（YAML）解析元数据：`title`、`servings`、`cooking_time`、`difficulty`、`tags`、`ingredients[]`、`steps[]`、`notes`
  - 内容分块策略：按章节/标题切分，生成 `schema.Document`（`Content` + `Metadata`）
- 向量数据库集成（Eino）
  - Embedding：OpenAI Embedding 组件（默认 `text-embedding-3-small`，可配）
  - Indexer/Retriever：首选 ES8（`index`、`topK`、`scoreThreshold` 可配）；后续可替换 VikingDB
- DeepSeek 模型集成（Eino OpenAI 模型）
  - `BaseURL=https://api.deepseek.com`，`model=deepseek-chat` 或 `deepseek-reasoner`（可配）
  - 问答流水线：检索 TopK → Prompt 构建 → ChatModel 推理 → 输出
  - 对话历史：内存会话（SessionID）与持久化接口预留

## REST API 原型（占位）
- `POST /api/v1/query`
  - 请求：`{ session_id?: string, query: string, top_k?: number }`
  - 响应：`{ answer: string, sources: [{id,score,metadata}], session_id }`
- `GET /api/v1/recipes`
  - 响应：`[{ id, title, tags, difficulty, cooking_time }]`

## 配置与安全
- 配置文件不含密钥；密钥均由环境变量注入
  - `DEEPSEEK_API_KEY`、`OPENAI_API_KEY`、`ES_USERNAME`、`ES_PASSWORD`
- 运行时仅以必要信息记录日志，不输出敏感字段

## 测试与性能基线
- 单元测试：
  - Parser：Front Matter/分块/容错（空字段、非法格式）
- 集成测试：
  - 替换 Embedding/Model 为 Mock，实现端到端（索引→检索→问答）
- 性能测试：
  - 基准：`POST /api/v1/query` 响应时间目标 `<500ms`（本地检索 + 流水线最小化），外部模型调用可能超出；支持流式/缓存优化

## 部署与交付
- Docker：多阶段构建，非 root 运行
- Kubernetes：`Deployment` + `Service` 模板（资源与环境变量）
- CI/CD：Go 测试 + Lint + Swagger 生成 + 镜像构建与推送

## 首次实现的交付物
- 完整目录与 go.mod
- 可启动的 `serve`（`/health` OK）与 `index` CLI 占位
- Viper 加载配置与环境变量覆盖
- Swagger 接入与初版 API Skeleton

## 后续迭代（第二阶段）
- 完成 RAG 问答链路、对话历史、检索调优
- 加入示例食谱模板与更完善的解析规则

参考资料：
- Eino GitHub 与文档：[1] https://github.com/cloudwego/eino；[2] https://www.cloudwego.io/docs/eino/core_modules/components/indexer_guide/；[3] https://www.cloudwego.io/docs/eino/ecosystem_integration/retriever/retriever_es8/；[4] https://www.cloudwego.io/docs/eino/core_modules/components/retriever_guide/；[5] https://www.cloudwego.io/docs/eino/
- DeepSeek API（OpenAI 兼容 BaseURL）：[6] https://api-docs.deepseek.com/