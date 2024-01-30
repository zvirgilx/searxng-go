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
}

var _engines = map[string]map[string]Engine{}

// RegisterEngine registers a search engine for used.
func RegisterEngine(name string, engine Engine, category string) {
	if _engines[category] == nil {
		_engines[category] = map[string]Engine{}
	}
	_engines[category][name] = engine
}

// GetEnginesByCategory gets an enable engines about a certain category.
func GetEnginesByCategory(category string) map[string]Engine {
	if es, ok := _engines[category]; ok {
		return es
	}
	return nil
}

// DisableEngine disables an engine.
func DisableEngine(category string, name string) {
	ces := _engines[category]
	if ces == nil {
		return
	}
	delete(ces, name)
}
