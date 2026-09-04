package daemon

import (
	"context"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting"
	"github.com/urfave/cli/v3"
)

func CmdConfig() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: i18n.Lang.T("cmd.daemon.config.Usage"),
		Commands: []*cli.Command{
			{
				Name:        "test",
				Usage:       i18n.Lang.T("cmd.daemon.config.test.Usage"),
				Description: i18n.Lang.T("cmd.daemon.config.test.Description"),
				Action:      CmdTestConfig,
			},
		},
	}
}

func CmdTestConfig(_ context.Context, _ *cli.Command) error {
	falseLogger := log.NewFalseLogger()

	testMain := testMainConfig()
	testDocker, dockerSupport := testDockerConfig()
	testFirewall := testFirewallConfig(dockerSupport)
	testAnalyzer := testAnalyzerConfig()
	testNotifications := testNotificationsConfig()
	testBlocklists := testBlocklistsConfig(falseLogger)
	testGeoip := testGeoipConfig(falseLogger)

	fmt.Println(
		"***\n"+i18n.Lang.T("cmd.daemon.config.test.settingTitle"),
		"\n", testMain,
		"\n", testFirewall,
		"\n", testAnalyzer,
		"\n", testNotifications,
		"\n", testDocker,
		"\n", testBlocklists,
		"\n", testGeoip,
		"\n***",
	)
	return nil
}

func testMainConfig() string {
	configTitle := i18n.Lang.T("cmd.daemon.config.test.main")
	if err := setting.Config.Validate(); err != nil {
		return resultError(configTitle, err)
	}

	if err := setting.Config.ValidateBeforeStart(); err != nil {
		return resultError(configTitle, err)
	}

	return resultOk(configTitle)
}

func testDockerConfig() (message string, dockerSupport bool) {
	configTitle := "docker"
	_, dockerSupport, err := setting.Config.OtherSettingsPath.ToDockerConfig(setting.Config.BinaryLocations)
	if err != nil {
		return resultError(configTitle, err), false
	}

	return resultOk(configTitle), dockerSupport
}

func testFirewallConfig(dockerSupport bool) string {
	configTitle := "firewall"
	if _, _, err := setting.Config.OtherSettingsPath.ToFirewallConfig(dockerSupport); err != nil {
		return resultError(configTitle, err)
	}

	return resultOk(configTitle)
}

func testAnalyzerConfig() string {
	configTitle := "analyzer"
	if _, err := setting.Config.OtherSettingsPath.ToAnalyzerConfig(setting.Config.BinaryLocations); err != nil {
		return resultError(configTitle, err)
	}

	return resultOk(configTitle)
}

func testNotificationsConfig() string {
	configTitle := "notifications"
	if _, err := setting.Config.OtherSettingsPath.ToNotificationsConfig(); err != nil {
		return resultError(configTitle, err)
	}

	return resultOk(configTitle)
}

func testBlocklistsConfig(logger log.Logger) string {
	configTitle := "blocklists"
	if _, _, err := setting.Config.OtherSettingsPath.ToBlocklistConfig(logger); err != nil {
		return resultError(configTitle, err)
	}

	return resultOk(configTitle)
}

func testGeoipConfig(logger log.Logger) string {
	configTitle := "geoip"
	if _, _, err := setting.Config.OtherSettingsPath.ToConfig(setting.Config.DataDir, logger); err != nil {
		return resultError(configTitle, err)
	}

	return resultOk(configTitle)
}

func resultOk(title string) string {
	return fmt.Sprintf("%s: \033[32mOk\033[0m", title)
}

func resultError(title string, err error) string {
	errText := i18n.Lang.T("cmd.daemon.config.test.error", map[string]interface{}{"Error": err})
	return fmt.Sprintf("%s: \033[31mError\n%s\u001B[0m", title, errText)
}
