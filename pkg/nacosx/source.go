package nacosx

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	nacosConfig "github.com/go-kratos/kratos/contrib/config/nacos/v2"
)

type Config struct {
	Addr      string
	Port      uint64
	Namespace string
	Env       string
	DataIds   []string // 支持多个 dataId
}

type Option func(*Config)

func WithAddr(addr string) Option {
	return func(c *Config) {
		c.Addr = addr
	}
}

func WithPort(port uint64) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithNamespace(namespace string) Option {
	return func(c *Config) {
		c.Namespace = namespace
	}
}

func WithEnv(env string) Option {
	return func(c *Config) {
		c.Env = env
	}
}

// WithDataId 设置单个 dataId
func WithDataId(dataId string) Option {
	return func(c *Config) {
		c.DataIds = []string{dataId}
	}
}

// WithDataIds 设置多个 dataId
func WithDataIds(dataIds ...string) Option {
	return func(c *Config) {
		c.DataIds = dataIds
	}
}

// NewNacosConfigSource 创建 Nacos 配置源（支持多个 dataId）
// 如果指定了多个 dataId，返回的第一个 Source 对应第一个 dataId，以此类推
func NewNacosConfigSource(opts ...Option) ([]config.Source, error) {
	cfg := &Config{
		Addr:      "127.0.0.1",
		Port:      8848,
		Namespace: "public",
		Env:       "dev",
		DataIds:   []string{"config.yaml"}, // 默认单个配置文件
	}
	for _, opt := range opts {
		opt(cfg)
	}

	cc := constant.ClientConfig{
		NamespaceId:         cfg.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "error",
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(cfg.Addr, cfg.Port),
	}

	// 创建配置客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, err
	}

	sources := make([]config.Source, 0, len(cfg.DataIds))
	for _, dataId := range cfg.DataIds {
		source := nacosConfig.NewConfigSource(client,
			nacosConfig.WithGroup(cfg.Env),
			nacosConfig.WithDataID(dataId))
		sources = append(sources, source)
	}

	return sources, nil
}
