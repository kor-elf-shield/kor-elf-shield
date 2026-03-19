package sources

import (
	"git.kor-elf.net/kor-elf-shield/blocklist"
	"git.kor-elf.net/kor-elf-shield/blocklist/parser"
)

type BlocklistSource interface {
	Get() (ipV4 parser.IPs, ipV6 parser.IPs, err error)
}

type blocklistSource struct {
	url    string
	parser parser.Parser
	config blocklist.Config
}

func NewBlocklistSource(url string, parser parser.Parser, config blocklist.Config) BlocklistSource {
	return &blocklistSource{
		url:    url,
		parser: parser,
		config: config,
	}
}

func (b *blocklistSource) Get() (ipV4 parser.IPs, ipV6 parser.IPs, err error) {
	return blocklist.GetSeparatedIPs(b.url, b.parser, b.config)
}

type blocklistSourceZip struct {
	url    string
	parser parser.Parser
	config blocklist.ConfigZip
}

func NewBlocklistSourceZip(url string, parser parser.Parser, config blocklist.ConfigZip) BlocklistSource {
	return &blocklistSourceZip{
		url:    url,
		parser: parser,
		config: config,
	}
}

func (b *blocklistSourceZip) Get() (ipV4 parser.IPs, ipV6 parser.IPs, err error) {
	return blocklist.GetZipSeparatedIPs(b.url, b.parser, b.config)
}
