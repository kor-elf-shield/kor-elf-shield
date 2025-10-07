package setting

type settingContract interface {
	GetLanguage() string
	GetFallbackLanguage() string
}

type setting struct {
	Language         string `mapstructure:"language"`
	FallbackLanguage string `mapstructure:"fallback_language"`
}

func (s *setting) GetLanguage() string {
	return s.Language
}

func (s *setting) GetFallbackLanguage() string {
	return s.FallbackLanguage
}
