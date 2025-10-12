package setting

type setting struct {
	Language         string `mapstructure:"language"`
	FallbackLanguage string `mapstructure:"fallback_language"`

	Log *log
}

func settingDefault() *setting {
	return &setting{
		Language:         "ru",
		FallbackLanguage: "ru",

		Log: logDefault(),
	}
}
