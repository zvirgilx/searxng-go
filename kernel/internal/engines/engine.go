package engines

import (
	"log/slog"

	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
)

func InitConfiguration(configuration map[string]map[string]engine.Config) {
	configuredEngines := map[string]map[string]engine.Engine{}

	for category, configMap := range configuration {
		engines := engine.GetEnginesByCategory(category)
		if len(engines) == 0 {
			continue
		}
		for name, conf := range configMap {
			if !conf.Enable {
				continue
			}
			if e, ok := engines[name]; ok {
				if err := e.ApplyConfig(conf); err != nil {
					slog.Error("failed to init configuration", slog.String("engineName", name), slog.String("error", err.Error()))
					continue
				}
				engine.RegisterTo(configuredEngines, e, category)
			}
		}
	}

	engine.SetGlobalEngines(configuredEngines)
}
