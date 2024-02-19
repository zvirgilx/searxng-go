package engine

import (
	"context"

	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

const (
	// CategoryGeneral search for general result, like direct search in google or bing.
	CategoryGeneral = "general"

	// CategoryVideo search for video result.
	CategoryVideo = "video"
)

type Engine interface {
	// Request reports how the engine initiates a request.
	Request(context.Context, *Options) error
	// Response reports how the engine parse the response.
	Response(context.Context, *Options, []byte) (*result.Result, error)

	GetName() string

	ApplyConfig(config Config) error
}

var _engines = map[string]map[string]Engine{}

// RegisterGlobalEngine registers a search engine for used.
func RegisterGlobalEngine(engine Engine, category string) {
	RegisterTo(_engines, engine, category)
}

func RegisterTo(engines map[string]map[string]Engine, engine Engine, category string) map[string]map[string]Engine {
	if engines[category] == nil {
		engines[category] = map[string]Engine{}
	}
	engines[category][engine.GetName()] = engine
	return engines
}

// GetEnginesByCategory gets an enable engines about a certain category.
func GetEnginesByCategory(category string) map[string]Engine {
	if es, ok := _engines[category]; ok {
		return es
	}
	return nil
}

func SetGlobalEngines(engines map[string]map[string]Engine) {
	_engines = engines
}
