package traits

import (
	_ "embed"
	"encoding/json"
)

//go:embed engine_traits.json
var traitsJson []byte

const LocaleAll = "all"

// EngineTraits is the engine trait. Most traits are languages and region.
// Traits all store in file `engine_traits.json`.
type EngineTraits struct {
	AllLocale string                       `json:"all_locale"` // AllLocale usually used as the default locale.
	Languages map[string]string            `json:"languages"`  // Languages contains engine language maps, e.g. "en"->"lang_en".
	Regions   map[string]string            `json:"regions"`    // Regions contains engine region maps, e.g. "en-US"->"US".
	Custom    map[string]map[string]string `json:"custom"`     // Custom contains engine custom information.
}

var traitsMap map[string]*EngineTraits

func InitTraits() error {
	var ts map[string]*EngineTraits
	if err := json.Unmarshal(traitsJson, &ts); err != nil {
		return err
	}
	traitsMap = ts
	return nil
}

func GetTrait(engineName string) *EngineTraits {
	return traitsMap[engineName]
}

func (e *EngineTraits) GetLanguage(locale string) string {
	if locale == LocaleAll && e.AllLocale != "" {
		return e.AllLocale
	}
	return e.Languages[locale]
}

func (e *EngineTraits) GetRegion(locale string) string {
	if locale == LocaleAll && e.AllLocale != "" {
		return e.AllLocale
	}
	return e.Regions[locale]
}

func (e *EngineTraits) GetCustom(field string) map[string]string {
	return e.Custom[field]
}
