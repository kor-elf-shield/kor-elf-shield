package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/cmd"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting"

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
		os.Exit(1)
	}

	err = i18n.InitLang(setting.Config.FallbackLanguage)
	if err != nil {
		fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		os.Exit(1)
	}

	err = i18n.Lang.ChangeLang(setting.Config.Language)
	if err != nil {
		fmt.Printf("\033[31m%s\033[0m\n", err.Error())
		os.Exit(1)
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
