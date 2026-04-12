package geoip

import (
	"errors"
	"fmt"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/geoip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"

	"git.kor-elf.net/kor-elf-shield/geoip2/service/maxmind/mmdb"
)

type Maxmind struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Interval int    `mapstructure:"interval"`
	URL      string `mapstructure:"url"`
	Language string `mapstructure:"language"`
}

func (m *Maxmind) ToConfig(dataDir string, log log.Logger) (*geoip.Config, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}

	url := m.URL
	if url == "" {
		url = mmdb.DownloadURLCityLite
	}
	download, err := mmdb.NewDownload(url, m.Username, m.Password, mmdb.DefaultDownloadConfig())
	if err != nil {
		return nil, err
	}

	dir := dataDir + "/geoip"

	language, err := m.language()
	if err != nil {
		return nil, err
	}

	logger := geoip.NewLogger(log)

	return &geoip.Config{
		GeoIP:    mmdb.NewCity(download, logger, dir, language),
		Interval: time.Duration(m.Interval) * time.Second,
	}, nil
}

func (m *Maxmind) validate() error {
	if m.Interval <= 0 {
		return errors.New("interval must be greater than 0")
	}

	if m.Username == "" {
		return errors.New("username is required")
	}

	if m.Password == "" {
		return errors.New("password is required")
	}

	if m.Language == "" {
		return errors.New("language is required")
	}

	return nil
}

func (m *Maxmind) language() (mmdb.Language, error) {
	switch m.Language {
	case "Russian":
		return mmdb.LanguageRussian, nil
	case "English":
		return mmdb.LanguageEnglish, nil
	case "Simplified Chinese":
		return mmdb.LanguageSimplifiedChinese, nil
	case "Brazilian Portuguese":
		return mmdb.LanguageBrazilianPortuguese, nil
	case "Spanish":
		return mmdb.LanguageSpanish, nil
	case "French":
		return mmdb.LanguageFrench, nil
	case "German":
		return mmdb.LanguageGerman, nil
	case "Japanese":
		return mmdb.LanguageJapanese, nil
	default:
		return mmdb.LanguageRussian, fmt.Errorf("invalid language: %s", m.Language)
	}
}
