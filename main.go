package main

import (
	"fmt"
	"kor-elf-shield/internal/cmd"
	"kor-elf-shield/internal/i18n"
	"kor-elf-shield/internal/setting"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// these flags will be set by the build flags
var (
	Version     = "development" // program version for this build
	MakeVersion = ""            // "make" program version if built with make
	ConfigPath  = "/etc/kor-elf-shield/kor-elf-shield.conf"
)

func init() {
	// Shell umask(0177)
	unix.Umask(0o177)

	setting.AppVer = Version
	setting.AppBuiltWith = formatBuiltWith()
	setting.AppStartTime = time.Now().UTC()
}

func main() {
	var err error

	configPath := findConfigPath(os.Args[1:])
	err = setting.InitSetting(configPath)
	if err != nil {
		fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		return
	}

	err = i18n.InitLang(setting.Config.GetFallbackLanguage())
	if err != nil {
		fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		return
	}

	err = i18n.Lang.ChangeLang(setting.Config.GetLanguage())
	if err != nil {
		fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		return
	}

	app := cmd.NewMainApp(cmd.AppVersion{Version: Version, Extra: formatBuiltWith()}, ConfigPath)
	_ = cmd.RunMainApp(app, os.Args...) // all errors should have been handled by the RunMainApp
}

func formatBuiltWith() string {
	version := runtime.Version()
	if len(MakeVersion) > 0 {
		version = MakeVersion + ", " + runtime.Version()
	}

	return " built with " + version
}

func findConfigPath(args []string) string {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--config=") {
			return strings.TrimPrefix(arg, "--config=")
		}
	}
	return ConfigPath
}
