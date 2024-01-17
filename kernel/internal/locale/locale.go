package locale

import (
	"log/slog"
	"strings"

	"github.com/zvirgilx/searxng-go/kernel/internal/engines/traits"
	"golang.org/x/text/language"
)

// GetLanguageFromTrait get language from locale and engine trait.
// locale is the input by user or get from browser header.
func GetLanguageFromTrait(locale string, trait *traits.EngineTraits, defaultVal string) string {
	log := slog.With("func", "GetLanguageFromTrait")
	if lang := trait.GetLanguage(locale); lang != "" {
		return lang
	}
	tag, err := language.Parse(locale)
	if err != nil {
		log.Error("msg", err)
		return defaultVal
	}

	base, _ := tag.Base()
	lang := base.String()

	// e.g. "zh-HK --> zh_Hant" or "zh-CN --> zh_Hans"
	if script, confidence := tag.Script(); confidence == language.Exact || confidence == language.Low {
		lang = base.String() + "_" + script.String()
	}
	if l := trait.GetLanguage(lang); l != "" {
		return l
	}
	return defaultVal
}

func ParseAcceptLanguage(acptLang string, defaultVal string) string {
	if acptLang == "*" || acptLang == "" {
		return defaultVal
	}
	lang := defaultVal

	// Multiple types, weighted with the quality value syntax:
	// Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5
	// Get ["fr-CH", " fr;q=0.9", ...]
	langQStrs := strings.Split(acptLang, ",")

	for _, langQStr := range langQStrs {
		trimedLangQStr := strings.Trim(langQStr, " ")

		// [fr,q=0.9], only need fr
		langQ := strings.Split(trimedLangQStr, ";")
		if langQ[0] != "*" {
			lang = langQ[0]
			break
		}
	}

	return lang
}
