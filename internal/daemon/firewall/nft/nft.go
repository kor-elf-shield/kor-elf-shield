package nft

import (
	nftables "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
)

type NFT interface {
	NewBuildBatch() (nft.BatchBuilder, error)
	RunBatch(batchBuilder nft.BatchBuilder) error
	RunBatchAndMoveFile(batchBuilder nft.BatchBuilder, fileDir string) error
	NFT() nft.NFT
}

type nftImpl struct {
	nft    nft.NFT
	tmpDir string
}

func New(nft nft.NFT, tmpDir string) NFT {
	return &nftImpl{
		nft:    nft,
		tmpDir: tmpDir,
	}
}

func (n *nftImpl) NFT() nft.NFT {
	return n.nft
}

func (n *nftImpl) NewBuildBatch() (nft.BatchBuilder, error) {
	return nftables.NewBatchBuilder(n.tmpDir)
}

func (n *nftImpl) RunBatch(batchBuilder nft.BatchBuilder) error {
	batch := batchBuilder.Build()
	defer func() {
		_ = batch.Close()
	}()

	return n.nft.ExecuteBatchAfterCheck(batch)
}

func (n *nftImpl) RunBatchAndMoveFile(batchBuilder nft.BatchBuilder, fileDir string) error {
	batch := batchBuilder.Build()
	defer func() {
		_ = batch.Close()
	}()
	if err := n.nft.ExecuteBatchAfterCheck(batch); err != nil {
		return err
	}

	return batch.MoveFile(fileDir)
}
