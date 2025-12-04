// 文件功能：HTTP 服务启动与路由配置；提供健康检查与基础 API 路由。
// 包功能：server 包，封装 CLI 与 HTTP 服务相关逻辑。
package server

import (
	"cook/internal/recipe/config"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// startHTTP：启动 HTTP 服务。
// 功能说明：加载应用配置，初始化 chi 路由与基础中间件，注册健康检查与占位 API。
// 参数说明：无。
// 返回值说明：
//   - error：监听失败时返回错误。
func startHTTP() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Get("/api/v1/recipes", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	})

	r.Post("/api/v1/query", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"answer":"TODO","sources":[]}`))
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	return http.ListenAndServe(addr, r)
}
