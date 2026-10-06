package validate

import "testing"

func TestValidate_PathFile_Success(t *testing.T) {
	if err := PathFile("/file", "test"); err != nil {
		t.Errorf("PathFile() should not return an error for valid path. %v", err)
	}
}

func TestValidate_PathFile_Error(t *testing.T) {
	invalidPathFile := []string{
		"invalid",
		"",
		" ",
		"./test",
		"../test",
	}

	for _, path := range invalidPathFile {
		t.Run(path, func(t *testing.T) {
			if err := PathFile(path, "test"); err == nil {
				t.Errorf("PathFile(%q) should return an error", path)
			}
		})
	}
}

func TestValidate_IsTomlFile_Success(t *testing.T) {
	validToml := []string{
		"/file.toml",
		"/file.Toml",
		"/file.TOML",
	}

	for _, path := range validToml {
		t.Run(path, func(t *testing.T) {
			if err := IsTomlFile(path, "test"); err != nil {
				t.Errorf("IsTomlFile(%q) should not return an error, got %v", path, err)
			}
		})
	}
}

func TestValidate_IsTomlFile_Error(t *testing.T) {
	invalidToml := []string{
		"file.toml",
		"/file",
	}

	for _, path := range invalidToml {
		t.Run(path, func(t *testing.T) {
			if err := IsTomlFile(path, "test"); err == nil {
				t.Errorf("IsTomlFile(%q) should return an error", path)
			}
		})
	}
}

func TestValidate_NftLimitRate_Success(t *testing.T) {
	validRates := []string{
		"1/s",
		"10/minute",
		"100/hour",
		"1000/day",
		"5/week",
		"0.5/second",
		"over 1/s",
		"10 packets/second",
		"1 kbytes/min",
		"2 mbytes/h",
		"1/s burst 5",
		"1/s burst 10 packets",
		"2/min burst 1.5 mbytes",
	}

	for _, rate := range validRates {
		t.Run(rate, func(t *testing.T) {
			if err := NftLimitRate(rate, "test"); err != nil {
				t.Errorf("NftLimitRate(%q) should not return an error, got: %v", rate, err)
			}
		})
	}
}

func TestValidate_NftLimitRate_Error(t *testing.T) {
	invalidRates := []string{
		"",
		"1",
		"/s",
		"1/",
		"one/s",
		"1/sec",
		"1 per second",
		"-1/s",
		"1/month",
		"1/s burst",
		"1/s burst x",
		"1/s burst 1 gb",
	}

	for _, rate := range invalidRates {
		t.Run(rate, func(t *testing.T) {
			if err := NftLimitRate(rate, "test"); err == nil {
				t.Errorf("NftLimitRate(%q) should return an error", rate)
			}
		})
	}
}

func TestValidate_Name_Success(t *testing.T) {
	validNames := []string{
		"a",
		"A",
		"abc",
		"name_1",
		"name-1",
		"Name_1-2",
		"abcdefghijklmnopqrstuvwxyzABCDEF", // 32 symbols
	}

	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			if err := Name(name, "test"); err != nil {
				t.Errorf("Name(%q) should not return an error, got %v", name, err)
			}
		})
	}
}

func TestValidate_Name_Error(t *testing.T) {
	invalidNames := []string{
		"",
		" ",
		"name.with.dot",
		"name/with/slash",
		"name with space",
		"name@with@at",
		"привет",                            // Cyrillic
		"abcdefghijklmnopqrstuvwxyzABCDEFG", // 33 symbols
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			if err := Name(name, "test"); err == nil {
				t.Errorf("Name(%q) should return an error", name)
			}
		})
	}
}

func TestValidate_NameWithDot_Success(t *testing.T) {
	validNames := []string{
		"a",
		"A",
		"abc",
		"name_1",
		"name-1",
		"name.with.dot",
		"Name_1-2.3",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // 64 symbols
	}

	for _, name := range validNames {
		t.Run(name, func(t *testing.T) {
			if err := NameWithDot(name, "test"); err != nil {
				t.Errorf("NameWithDot(%q) should not return an error, got %v", name, err)
			}
		})
	}
}

func TestValidate_NameWithDot_Error(t *testing.T) {
	invalidNames := []string{
		"",
		" ",
		"name/with/slash",
		"name with space",
		"name@with@at",
		"привет", // Cyrillic
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // 65 symbols
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			if err := NameWithDot(name, "test"); err == nil {
				t.Errorf("NameWithDot(%q) should return an error", name)
			}
		})
	}
}

func TestValidate_Port_Success(t *testing.T) {
	validPorts := []int{
		0,
		1,
		80,
		443,
		65535,
	}

	for _, port := range validPorts {
		t.Run("port", func(t *testing.T) {
			if err := Port(port, "test"); err != nil {
				t.Errorf("Port(%d) should not return an error, got %v", port, err)
			}
		})
	}
}

func TestValidate_Port_Error(t *testing.T) {
	invalidPorts := []int{
		-1,
		65536,
		70000,
	}

	for _, port := range invalidPorts {
		t.Run("port", func(t *testing.T) {
			if err := Port(port, "test"); err == nil {
				t.Errorf("Port(%d) should return an error", port)
			}
		})
	}
}
