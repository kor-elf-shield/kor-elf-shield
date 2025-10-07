package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

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
	ChangeLang(lang string) error
}

type lang struct {
	loadedLangs map[string]bool

	bundler           *i18n.Bundle
	localizer         *i18n.Localizer
	fallbackLocalizer *i18n.Localizer

	fallbackLang string
	currentLang  string
}

func InitLang(fallbackLang string) error {
	bundle := i18n.NewBundle(language.Make(fallbackLang))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	newLang := lang{
		loadedLangs: make(map[string]bool),
		bundler:     bundle,

		fallbackLang: fallbackLang,
		currentLang:  fallbackLang,
	}

	err := newLang.ensureLocaleLoaded(fallbackLang)
	if err != nil {
		return err
	}

	localizer := i18n.NewLocalizer(bundle, fallbackLang)
	newLang.localizer = localizer
	newLang.fallbackLocalizer = localizer

	Lang = &newLang
	return nil
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

func (l *lang) ChangeLang(lang string) error {
	if lang == l.currentLang {
		return nil
	}

	err := l.ensureLocaleLoaded(lang)
	if err != nil {
		return err
	}

	l.localizer = i18n.NewLocalizer(l.bundler, lang)
	l.currentLang = lang
	return nil
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

// We upload translation files as needed.
func (l *lang) ensureLocaleLoaded(lang string) error {
	if l.loadedLangs[lang] {
		return nil
	}

	path := fmt.Sprintf("locales/locale.%s.json", lang)
	_, err := l.bundler.LoadMessageFileFS(localeFiles, path)
	if err != nil {
		return err
	}

	l.loadedLangs[lang] = true
	return nil
}
