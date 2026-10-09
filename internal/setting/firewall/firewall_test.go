package firewall

import (
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
)

func TestFirewall_InitSetting_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "firewall.toml")

	if err := testutil.SaveTomlEmpty(path); err != nil {
		t.Fatalf("failed to create toml file: %v", err)
	}

	_, err := InitSetting(path)
	if err != nil {
		t.Errorf("InitSetting() should not return an error, got %v", err)
		return
	}
}

func TestFirewall_InitSetting_Error(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		if _, err := InitSetting("firewall.toml"); err == nil {
			t.Errorf("InitSetting() should return an error for non-absolute path")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "not_exists.toml")

		if _, err := InitSetting(path); err == nil {
			t.Errorf("InitSetting() should return an error for missing file")
		}
	})

	t.Run("file not TOML", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "not_toml.txt")
		if err := testutil.SaveTomlEmpty(path); err != nil {
			t.Errorf("testutil.SaveTomlEmpty() failed: %v", err)
		}

		if _, err := InitSetting(path); err == nil {
			t.Errorf("InitSetting() should return an error for non-TOML file")
		}
	})
}

func TestFirewall_Validate_Success(t *testing.T) {
	dir := t.TempDir()
	s := defaultFirewallForTest(dir)
	if err := s.Validate(); err != nil {
		t.Errorf("Setting.Validate() should not return an error, got %v", err)
	}
}

func TestFirewall_Validate_Error(t *testing.T) {
	t.Run("empty IP4", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.IP4 = ip4{}
		if err := s.Validate(); err == nil {
			t.Errorf("Setting.Validate() should return an error for empty IP4")
		}
	})

	t.Run("empty Options", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.Options = options{}
		if err := s.Validate(); err == nil {
			t.Errorf("Setting.Validate() should return an error for empty Options")
		}
	})

	t.Run("empty MetadataNaming", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.MetadataNaming = metadataNaming{}
		if err := s.Validate(); err == nil {
			t.Errorf("Setting.Validate() should return an error for empty MetadataNaming")
		}
	})

	t.Run("empty Policy", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.Policy = policy{}
		if err := s.Validate(); err == nil {
			t.Errorf("Setting.Validate() should return an error for empty Policy")
		}
	})

	t.Run("empty RulesGuard", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.RulesGuard = RulesGuard{}
		if err := s.Validate(); err == nil {
			t.Errorf("Setting.Validate() should return an error for empty RulesGuard")
		}
	})
}

func TestFirewall_ToPorts_Success(t *testing.T) {
	t.Run("empty Ports", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.Ports = []Port{}
		if _, _, err := s.ToPorts(); err != nil {
			t.Errorf("Setting.ToPorts() error = %v, wantErr nil", err)
		}
	})

	t.Run("not empty Ports", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.Ports = []Port{
			*defaultPortForTest(),
		}
		if _, _, err := s.ToPorts(); err != nil {
			t.Errorf("Setting.ToPorts() error = %v, wantErr nil", err)
		}
	})
}

func TestFirewall_ToPorts_Error(t *testing.T) {
	t.Run("empty Ports", func(t *testing.T) {
		port := *defaultPortForTest()
		port.Numbers = []int{}

		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.Ports = []Port{
			port,
		}
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Setting.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Ports", func(t *testing.T) {
		port := *defaultPortForTest()
		port.Numbers = []int{-1}

		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.Ports = []Port{
			port,
		}
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Setting.ToPorts() error = nil, wantErr not nil")
		}
	})
}

func TestFirewall_ToIPs_Success(t *testing.T) {
	t.Run("empty IPS", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.IPS = []IP{}
		if _, err := s.ToIPs(); err != nil {
			t.Errorf("Setting.ToIPs() error = %v, wantErr nil", err)
		}
	})

	t.Run("not empty IPS", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.IPS = []IP{
			*defaultIPForTest(),
		}
		if _, err := s.ToIPs(); err != nil {
			t.Errorf("Setting.ToIPs() error = %v, wantErr nil", err)
		}
	})
}

func TestFirewall_ToIPs_Error(t *testing.T) {
	t.Run("empty IP", func(t *testing.T) {
		ip := *defaultIPForTest()
		ip.IPs = []string{}

		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.IPS = []IP{
			ip,
		}
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("Setting.ToIPs() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid IP", func(t *testing.T) {
		ip := *defaultIPForTest()
		ip.IPs = []string{
			"invalid",
		}

		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.IPS = []IP{
			ip,
		}
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("Setting.ToIPs() error = nil, wantErr not nil")
		}
	})
}

func TestFirewall_ToConfigPortKnocking_Success(t *testing.T) {
	t.Run("empty IPS", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.PortKnocking = []portKnocking{}
		if _, err := s.ToConfigPortKnocking(); err != nil {
			t.Errorf("Setting.ToConfigPortKnocking() error = %v, wantErr nil", err)
		}
	})

	t.Run("not empty IPS", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.PortKnocking = []portKnocking{
			*defaultPortKnockingForTest(),
		}
		if _, err := s.ToConfigPortKnocking(); err != nil {
			t.Errorf("Setting.ToConfigPortKnocking() error = %v, wantErr nil", err)
		}
	})
}

func TestFirewall_ToConfigPortKnocking_Error(t *testing.T) {
	t.Run("empty Knocks", func(t *testing.T) {
		portKnockingTest := *defaultPortKnockingForTest()
		portKnockingTest.Knocks = []portKnockingKnock{}

		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.PortKnocking = []portKnocking{
			portKnockingTest,
		}
		if _, err := s.ToConfigPortKnocking(); err == nil {
			t.Errorf("Setting.ToConfigPortKnocking() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Knocks", func(t *testing.T) {
		portKnockingTest := *defaultPortKnockingForTest()
		portKnockingTest.Port = -1

		dir := t.TempDir()
		s := defaultFirewallForTest(dir)
		s.PortKnocking = []portKnocking{
			portKnockingTest,
		}
		if _, err := s.ToConfigPortKnocking(); err == nil {
			t.Errorf("Setting.ToConfigPortKnocking() error = nil, wantErr not nil")
		}
	})
}

func defaultFirewallForTest(dir string) *Setting {
	return &Setting{
		Ports:          []Port{},
		IPS:            []IP{},
		IP4:            *defaultIp4ForTest(),
		IP6:            *defaultIp6ForTest(),
		Options:        *defaultOptionsForTest(dir),
		MetadataNaming: *defaultMetadataNamingForTest(),
		Policy:         *defaultPolicyForTest(),
		PortKnocking:   []portKnocking{},
		RulesGuard:     *defaultRulesGuardForTest(),
	}
}
