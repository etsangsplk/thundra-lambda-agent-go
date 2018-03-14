package thundra

import (
	"thundra-agent-go/plugin"
	"os"
)

type Builder interface {
	AddPlugin(plugin.Plugin) Builder
	SetReporter(Reporter) Builder
	SetAPIKey(string) Builder
	Build() *thundra
}

type builder struct {
	plugins  []plugin.Plugin
	reporter Reporter
	apiKey   string
}

func (b *builder) AddPlugin(plugin plugin.Plugin) Builder {
	b.plugins = append(b.plugins, plugin)
	return b
}

func (b *builder) SetReporter(reporter Reporter) Builder {
	b.reporter = reporter
	return b
}

func (b *builder) SetAPIKey(apiKey string) Builder {
	b.apiKey = apiKey
	return b
}

func (b *builder) Build() *thundra {
	if b.reporter == nil {
		b.reporter = &reporterImpl{}
	}
	if b.apiKey == "" {
		b.apiKey = os.Getenv(thundraApiKey)
	}
	return &thundra{
		b.plugins,
		b.reporter,
		b.apiKey,
	}
}

func NewBuilder() Builder {
	return &builder{}
}
