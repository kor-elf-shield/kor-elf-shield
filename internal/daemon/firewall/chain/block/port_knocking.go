package block

import (
	"strconv"
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
)

func NewPortKnocking(nft nft.NFT, family family.Type, table string, name string, ipVersion ip.Version, timeout uint32) error {
	params := []string{"type", ipVersion.ToNftForSet() + ";", "flags timeout; timeout", strconv.Itoa(int(timeout)) + "s;"}
	_, err := newList(nft, family, table, name, strings.Join(params, " "))
	if err != nil {
		return err
	}

	return nil
}
