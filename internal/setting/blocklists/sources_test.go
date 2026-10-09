package blocklists

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestSources_ToSourceConfig_Success(t *testing.T) {
	s := defaultSourcesForTest()
	exclusionChecker := newExclusionCheckerMock()
	if _, err := s.ToSourceConfig(exclusionChecker); err != nil {
		t.Errorf("Sources.ToSourceConfig() should not return an error, got %v", err)
	}
}

func TestSources_ToSourceConfig_Error(t *testing.T) {
	invalidName := []string{
		"",
		" ",
		"name.with.dot",
		"name/with/slash",
		"name with space",
		"name@with@at",
		"привет",                            // Cyrillic
		"abcdefghijklmnopqrstuvwxyzABCDEFG", // 33 symbols
	}
	for _, name := range invalidName {
		t.Run(fmt.Sprintf("invalid name %s", name), func(t *testing.T) {
			s := defaultSourcesForTest()
			s.Name = name
			exclusionChecker := newExclusionCheckerMock()
			if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
				t.Errorf("Sources.ToSourceConfig() should return an error for name %s", name)
			}
		})
	}

	invalidUrl := []string{
		"",
		" ",
		"invalid url",
		"/test",
	}
	for _, url := range invalidUrl {
		t.Run(fmt.Sprintf("invalid url %s", url), func(t *testing.T) {
			s := defaultSourcesForTest()
			s.URL = url
			exclusionChecker := newExclusionCheckerMock()
			if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
				t.Errorf("Sources.ToSourceConfig() should return an error for url %s", url)
			}
		})
	}

	invalidLimit := []int{-1, -1000}
	for _, limit := range invalidLimit {
		t.Run(fmt.Sprintf("invalid limit %d", limit), func(t *testing.T) {
			s := defaultSourcesForTest()
			s.Limit = limit
			exclusionChecker := newExclusionCheckerMock()
			if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
				t.Errorf("Sources.ToSourceConfig() should return an error for limit %d", limit)
			}
		})
	}

	invalidInterval := []int64{59, 0, -1, -100}
	for _, interval := range invalidInterval {
		t.Run(fmt.Sprintf("invalid interval %d", interval), func(t *testing.T) {
			s := defaultSourcesForTest()
			s.Interval = interval
			exclusionChecker := newExclusionCheckerMock()
			if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
				t.Errorf("Sources.ToSourceConfig() should return an error for interval %d", interval)
			}
		})
	}

	invalidFormat := []string{
		"",
		" ",
		"invalid",
	}
	for _, format := range invalidFormat {
		t.Run(fmt.Sprintf("invalid format %s", format), func(t *testing.T) {
			s := defaultSourcesForTest()
			s.Format = format
			exclusionChecker := newExclusionCheckerMock()
			if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
				t.Errorf("Sources.ToSourceConfig() should return an error for format %s", format)
			}
		})
	}
}

func TestSources_ToSourceConfig_Error_Json(t *testing.T) {
	t.Run("empty JsonField", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "json"
		s.JsonField = ""
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})
}

func TestSources_ToSourceConfig_Error_Txt(t *testing.T) {
	t.Run("empty TxtSeparator", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "txt"
		s.TxtSeparator = ""
		s.TxtFieldIP = 0
		s.TxtType = "default"
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("negative TxtFieldIP", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "txt"
		s.TxtSeparator = "\t"
		s.TxtFieldIP = -1
		s.TxtType = "default"
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("empty TxtType", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "txt"
		s.TxtSeparator = "\t"
		s.TxtFieldIP = 0
		s.TxtType = ""
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid TxtType", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "txt"
		s.TxtSeparator = "\t"
		s.TxtFieldIP = 0
		s.TxtType = "invalid"
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("negative TxtFieldIp2", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "txt"
		s.TxtSeparator = "\t"
		s.TxtFieldIP = 0
		s.TxtType = "interval"
		s.TxtFieldIp2 = -1
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("negative TxtFieldCIDR", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "txt"
		s.TxtSeparator = "\t"
		s.TxtFieldIP = 0
		s.TxtType = "cidr"
		s.TxtFieldCIDR = -1
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})
}

func TestSources_ToSourceConfig_Error_Rss(t *testing.T) {
	t.Run("empty RssTag", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "rss"
		s.RssTag = ""
		s.RssField = "title"
		s.RssFieldSeparator = ""
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("empty RssField", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "rss"
		s.RssTag = "item"
		s.RssField = ""
		s.RssFieldSeparator = ""
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("negative RssFieldIP", func(t *testing.T) {
		s := defaultSourcesForTest()
		s.Format = "rss"
		s.RssTag = "item"
		s.RssField = "title"
		s.RssFieldSeparator = " "
		s.RssFieldIP = -1
		exclusionChecker := newExclusionCheckerMock()
		if _, err := s.ToSourceConfig(exclusionChecker); err == nil {
			t.Errorf("Sources.ToSourceConfig() error = nil, wantErr not nil")
		}
	})
}

func defaultSourcesForTest() *Sources {
	sources := defaultBaseSourcesForTest()
	sources.Format = fakeRandFormat()

	switch sources.Format {
	case "json":
		return fakeFormatJson(sources)
	case "txt":
		return fakeFormatTxt(sources)
	case "rss":
		return fakeFormatRss(sources)
	}

	return sources
}

func defaultBaseSourcesForTest() *Sources {
	interval := fakers.RandInt(60, 86400)
	return &Sources{
		Enabled:  true,
		Name:     "test",
		URL:      "https://test.localhost",
		Limit:    fakers.RandInt(0, 10000000),
		Interval: int64(interval),
		Zip:      fakers.RandBool(),
	}
}

func fakeRandFormat() string {
	items := []string{"json", "txt", "rss"}
	return fakers.RandItem(items)
}

func fakeFormatJson(s *Sources) *Sources {
	s.JsonField = "cidr"

	return s
}

func fakeFormatTxt(s *Sources) *Sources {
	s.TxtFieldIP = fakers.RandInt(0, 1000)
	s.TxtSeparator = "\t"

	s.TxtType = fakeRandTxtType()
	if s.TxtType == "interval" {
		s.TxtFieldIp2 = fakers.RandInt(0, 1000)
	} else if s.TxtType == "cidr" {
		s.TxtFieldCIDR = fakers.RandInt(0, 1000)
	}

	return s
}

func fakeRandTxtType() string {
	items := []string{"default", "cidr", "interval"}
	return fakers.RandItem(items)
}

func fakeFormatRss(s *Sources) *Sources {
	s.RssTag = "item"
	s.RssField = "title"
	s.RssFieldSeparator = fakeRandRssFieldSeparator()

	if s.RssFieldSeparator != "" {
		s.RssFieldIP = fakers.RandInt(0, 1000)
	}

	return s
}

func fakeRandRssFieldSeparator() string {
	items := []string{"", "|"}
	return fakers.RandItem(items)
}

type exclusionCheckerResult struct {
	excluded bool
	ips      []string
	err      error
}

type exclusionCheckerMock struct {
	byIP            map[string]exclusionCheckerResult
	defaultExcluded bool
	defaultIPs      []string
	defaultErr      error
	calls           []string
}

func (e *exclusionCheckerMock) IsExcluded(ip string) (excluded bool, ips []string, err error) {
	e.calls = append(e.calls, ip)

	if result, ok := e.byIP[ip]; ok {
		return result.excluded, append([]string(nil), result.ips...), result.err
	}

	return e.defaultExcluded, append([]string(nil), e.defaultIPs...), e.defaultErr
}

func newExclusionCheckerMock() *exclusionCheckerMock {
	return &exclusionCheckerMock{
		byIP:  make(map[string]exclusionCheckerResult),
		calls: make([]string, 0),
	}
}
