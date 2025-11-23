package config

import (
    "strings"
    "github.com/spf13/viper"
)

type ServerConfig struct {
    Port int `mapstructure:"port"`
}

type DeepSeekConfig struct {
    BaseURL string `mapstructure:"base_url"`
    Model   string `mapstructure:"model"`
    APIKey  string `mapstructure:"api_key"`
}

type ES8Config struct {
    Address  string `mapstructure:"address"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Index    string `mapstructure:"index"`
}

type AppConfig struct {
    Server   ServerConfig  `mapstructure:"server"`
    DeepSeek DeepSeekConfig `mapstructure:"deepseek"`
    ES8      ES8Config     `mapstructure:"es8"`
}

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

    if k := v.GetString("DEEPSEEK_API_KEY"); k != "" {
        cfg.DeepSeek.APIKey = k
    }
    return cfg, nil
}

func setDefaults(v *viper.Viper) error {
    v.SetDefault("server.port", 8080)
    v.SetDefault("deepseek.base_url", "https://api.deepseek.com")
    v.SetDefault("deepseek.model", "deepseek-chat")
    v.SetDefault("es8.index", "recipes")
    return nil
}