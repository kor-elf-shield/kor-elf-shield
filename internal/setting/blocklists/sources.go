package blocklists

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"git.kor-elf.net/kor-elf-shield/blocklist"
	"git.kor-elf.net/kor-elf-shield/blocklist/parser"
	daemonBlocklist "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/blocklist"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/blocklist/sources"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type Sources struct {
	Enabled  bool   `mapstructure:"enabled"`
	Name     string `mapstructure:"name"`
	URL      string `mapstructure:"url"`
	Limit    int    `mapstructure:"limit"`
	Interval int64  `mapstructure:"interval"`
	Zip      bool   `mapstructure:"zip"`
	Format   string `mapstructure:"format"`

	JsonField string `mapstructure:"json_field"`

	TxtType      string `mapstructure:"txt_type"`
	TxtFieldIP   int    `mapstructure:"txt_field_ip"`
	TxtFieldIp2  int    `mapstructure:"txt_field_ip2"`
	TxtFieldCIDR int    `mapstructure:"txt_field_cidr"`
	TxtSeparator string `mapstructure:"txt_separator"`

	RssTag            string `mapstructure:"rss_tag"`
	RssField          string `mapstructure:"rss_field"`
	RssFieldIP        int    `mapstructure:"rss_field_ip"`
	RssFieldSeparator string `mapstructure:"rss_field_separator"`
}

func (s *Sources) ToSourceConfig() (*daemonBlocklist.SourceConfig, error) {
	if err := s.Validate(); err != nil {
		return &daemonBlocklist.SourceConfig{}, err
	}

	pars, err := s.parser()
	if err != nil {
		return &daemonBlocklist.SourceConfig{}, err
	}

	config := blocklist.NewConfig(uint(s.Limit))
	if s.TxtType == "interval" {
		config.Validator = &parser.IPRangeValidator{}
	}

	if s.Zip {
		configZip := blocklist.NewConfigZip(config)
		return &daemonBlocklist.SourceConfig{
			Name:     s.Name,
			Interval: time.Duration(s.Interval) * time.Second,
			Source:   sources.NewBlocklistSourceZip(s.URL, pars, configZip),
		}, nil
	}

	return &daemonBlocklist.SourceConfig{
		Name:     s.Name,
		Interval: time.Duration(s.Interval) * time.Second,
		Source:   sources.NewBlocklistSource(s.URL, pars, config),
	}, nil
}

func (s *Sources) Validate() error {
	if err := validate.Name(s.Name, "sources.name"); err != nil {
		return err
	}

	if s.Interval < 60 {
		return errors.New("invalid limit. Must be greater than or equal to 60 seconds")
	}

	if s.Limit < 0 {
		return errors.New("invalid limit. Must be greater than or equal to 0")
	}

	if s.URL == "" {
		return errors.New("url is required")
	} else if !strings.HasPrefix(s.URL, "http://") && !strings.HasPrefix(s.URL, "https://") {
		return errors.New("the URL must be to an HTTP or HTTPS resource")
	}

	return nil
}

func (s *Sources) parser() (parser.Parser, error) {
	switch s.Format {
	case "json":
		if s.JsonField == "" {
			return nil, errors.New("json_field is required")
		}
		return parserJson(s.JsonField)
	case "txt":
		return s.parserText()
	case "rss":
		if s.RssTag == "" {
			return nil, errors.New("rss_tag is required")
		}
		if s.RssField == "" {
			return nil, errors.New("rss_field is required")
		}
		if s.RssFieldIP < 0 {
			return nil, errors.New("rss_field_ip must be greater than or equal to 0")
		}
		return parserRss(s.RssTag, s.RssField, uint(s.RssFieldIP), s.RssFieldSeparator)
	}

	return nil, fmt.Errorf("format not support")
}

func (s *Sources) parserText() (parser.Parser, error) {
	if s.TxtType == "" {
		return nil, errors.New("txt_type is required")
	}
	if s.TxtFieldIP < 0 {
		return nil, errors.New("txt_field_ip must be greater than or equal to 0")
	}
	if s.TxtSeparator == "" {
		return nil, errors.New("txt_separator is required")
	}
	switch s.TxtType {
	case "default":
		return parserTextDefault(uint8(s.TxtFieldIP), s.TxtSeparator)
	case "cidr":
		if s.TxtFieldCIDR < 0 {
			return nil, errors.New("txt_field_cidr must be greater than or equal to 0")
		}
		return parserTextCIDR(uint8(s.TxtFieldIP), uint8(s.TxtFieldCIDR), s.TxtSeparator)
	case "interval":
		if s.TxtFieldIp2 < 0 {
			return nil, errors.New("txt_field_ip2 must be greater than or equal to 0")
		}
		return parserTextInterval(uint8(s.TxtFieldIP), uint8(s.TxtFieldIp2), s.TxtSeparator)
	}
	return nil, fmt.Errorf("txt_type not support")
}

func parserJson(fieldName string) (parser.Parser, error) {
	return parser.NewJsonLines(func(item json.RawMessage) (string, error) {
		var line map[string]any
		if err := json.Unmarshal(item, &line); err != nil {
			return "", fmt.Errorf("unmarshal json item: %w", err)
		}

		v, ok := line[fieldName]
		if !ok {
			return "", nil
		}

		ip, ok := v.(string)
		if !ok {
			return "", nil
		}

		return ip, nil
	})
}

func parserTextDefault(ip uint8, separator string) (parser.Parser, error) {
	return parser.NewText(parser.NewDefaultTextExtract(ip, separator))
}

func parserTextCIDR(ip uint8, cidr uint8, separator string) (parser.Parser, error) {
	textExtract := parser.NewCIDRTextExtract(ip, cidr, separator)
	return parser.NewText(textExtract)
}

func parserTextInterval(ip uint8, ip2 uint8, separator string) (parser.Parser, error) {
	textExtract := parser.NewIntervalTextExtract(ip, ip2, separator)
	return parser.NewText(textExtract)
}

func parserRss(itemTag string, fieldName string, fieldIP uint, separator string) (parser.Parser, error) {
	return parser.NewRss(func(decoder *xml.Decoder, start xml.StartElement) (string, error) {
		for {
			tok, err := decoder.Token()
			if err != nil {
				if err == io.EOF {
					return "", nil
				}
				return "", err
			}

			switch t := tok.(type) {
			case xml.StartElement:
				if t.Name.Local != itemTag {
					continue
				}

				value, err := parserRssReadFieldFromItem(decoder, t, fieldName)
				if err != nil {
					return "", err
				}
				if value != "" {
					if separator != "" {
						fields := strings.Split(value, separator)
						if len(fields) <= int(fieldIP) {
							return "", nil
						}
						return strings.TrimSpace(fields[fieldIP]), nil
					}
					return strings.TrimSpace(value), nil
				}
			}
		}
	})
}

func parserRssReadFieldFromItem(decoder *xml.Decoder, start xml.StartElement, fieldTag string) (string, error) {
	depth := 1

	for {
		tok, err := decoder.Token()
		if err != nil {
			return "", err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			depth++

			if t.Name.Local == fieldTag {
				var value string
				if err := decoder.DecodeElement(&value, &t); err != nil {
					return "", err
				}
				return strings.TrimSpace(value), nil
			}

		case xml.EndElement:
			depth--
			if depth == 0 && t.Name.Local == start.Name.Local {
				return "", nil
			}
		}
	}
}
