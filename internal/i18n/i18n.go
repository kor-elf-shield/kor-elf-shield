package i18n

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	Lang langContract

	//go:embed locales/*.json
	localeFiles embed.FS
)

type langContract interface {
	T(messageID string, data ...map[string]interface{}) string
	ChangeLang(lang string)
}

type lang struct {
	bundler           *i18n.Bundle
	localizer         *i18n.Localizer
	fallbackLocalizer *i18n.Localizer

	fallbackLang string
	currentLang  string
}

func InitLang(fallbackLang string) {
	bundle := i18n.NewBundle(language.Make(fallbackLang))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, _ = bundle.LoadMessageFileFS(localeFiles, "locales/locale.ru.json")
	_, _ = bundle.LoadMessageFileFS(localeFiles, "locales/locale.en.json")

	localizer := i18n.NewLocalizer(bundle, fallbackLang)

	Lang = &lang{
		bundler:           bundle,
		localizer:         localizer,
		fallbackLocalizer: localizer,

		fallbackLang: fallbackLang,
		currentLang:  fallbackLang,
	}
}

func (l *lang) T(messageID string, data ...map[string]interface{}) string {
	cfg := &i18n.LocalizeConfig{
		MessageID: messageID,
	}
	if len(data) > 0 {
		cfg.TemplateData = data[0]
	}

	msg, err := l.localizer.Localize(cfg)
	if err != nil {
		return l.fallbackT(messageID, data...)
	}
	return msg
}

func (l *lang) ChangeLang(lang string) {
	if lang == l.currentLang {
		return
	}

	l.localizer = i18n.NewLocalizer(l.bundler, lang)
	l.currentLang = lang
}

func (l *lang) fallbackT(messageID string, data ...map[string]interface{}) string {
	cfg := &i18n.LocalizeConfig{
		MessageID: messageID,
	}
	if len(data) > 0 {
		cfg.TemplateData = data[0]
	}

	msg, err := l.fallbackLocalizer.Localize(cfg)
	if err != nil {
		return messageID
	}
	return msg
}
