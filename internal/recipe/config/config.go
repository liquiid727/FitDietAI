// 文件功能：应用配置加载与默认值设置；支持从配置文件与环境变量合并生成运行时配置。
// 用户认证相关处理
// 包功能：配置包，提供 Server、DeepSeek、ES8 等模块的配置结构体与加载函数。
package config

import (
	"strings"

	"github.com/spf13/viper"
)

// ServerConfig：HTTP 服务配置。
//   - Port：服务监听端口。
type ServerConfig struct {
	Port int `mapstructure:"port"` // 服务监听端口
}

// DeepSeekConfig：DeepSeek 大模型访问配置。
//   - BaseURL：服务基础地址；
//   - Model：模型名称；
//   - APIKey：访问密钥，可由环境变量覆盖。
type DeepSeekConfig struct {
	BaseURL string `mapstructure:"base_url"` // DeepSeek 接口基础地址
	Model   string `mapstructure:"model"`    // 模型名称
	APIKey  string `mapstructure:"api_key"`  // 访问密钥
}

// ES8Config：Elasticsearch v8 连接配置。
//   - Address：服务地址；
//   - Username：用户名；
//   - Password：密码；
//   - Index：默认索引名称。
type ES8Config struct {
	Address  string `mapstructure:"address"`  // ES 地址
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
	Index    string `mapstructure:"index"`    // 索引名称
}

// AppConfig：应用配置根结构。
//   - Server：HTTP 服务配置；
//   - DeepSeek：大模型调用配置；
//   - ES8：向量检索/索引构建的存储后端配置。
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`   // 服务配置
	DeepSeek DeepSeekConfig `mapstructure:"deepseek"` // DeepSeek 配置
	ES8      ES8Config      `mapstructure:"es8"`      // ES8 配置
}

// Load：加载应用配置。
// 功能说明：从 configs/config.yaml 与环境变量（前缀 RECIPE_AGENT）读取配置，应用默认值与环境覆盖，返回统一结构。
// 参数说明：无。
// 返回值说明：
//   - *AppConfig：完整应用配置；
//   - error：读取或反序列化失败时返回错误。
func Load() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("configs")
	v.SetConfigType("yaml")

	v.SetEnvPrefix("RECIPE_AGENT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = setDefaults(v)
	_ = v.ReadInConfig()

	cfg := new(AppConfig)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// 环境变量覆盖：优先使用 RECIPE_AGENT_DEEPSEEK_API_KEY
	if k := v.GetString("DEEPSEEK_API_KEY"); k != "" {
		cfg.DeepSeek.APIKey = k
	}
	return cfg, nil
}

// setDefaults：设置默认配置值。
// 功能说明：为未提供的配置项填充合理默认值；便于本地快速运行。
// 参数说明：
//   - v：viper 实例。
//
// 返回值说明：
//   - error：当前实现不返回错误，预留扩展。
func setDefaults(v *viper.Viper) error {
	v.SetDefault("server.port", 8080)
	v.SetDefault("deepseek.base_url", "https://api.deepseek.com")
	v.SetDefault("deepseek.model", "deepseek-chat")
	v.SetDefault("es8.index", "recipes")
	return nil
}
