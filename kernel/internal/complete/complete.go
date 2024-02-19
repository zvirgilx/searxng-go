package complete

import (
	"context"
)

const (
	TypeText  = "text"
	TypeMedia = "media"
)

type Config struct {
	EnableEngines []string `mapstructure:"enable_engines"`
}

// Completer defines an engine that completes a search query.
type Completer interface {
	// Complete returns the completed results according to the search query and locale.
	Complete(ctx context.Context, query string, locale string) []Result
}

var completers = map[string]Completer{}

func RegisterCompleter(name string, completer Completer) {
	completers[name] = completer
}

func InitCompleters(conf Config) {
	enable := make(map[string]Completer)
	for _, name := range conf.EnableEngines {
		if f, ok := completers[name]; ok {
			enable[name] = f
		}
	}
	completers = enable
}

func Complete(ctx context.Context, q string, locale string) []Result {
	var results []Result
	for ci := range completers {
		res := completers[ci].Complete(ctx, q, locale)
		results = mergeResult(results, res)
	}
	return results
}

// TODO: optimize the result merge
func mergeResult(src []Result, tar []Result) []Result {
	return append(src, tar...)
}
